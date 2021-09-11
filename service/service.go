package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net"
	"os"
	"os/signal"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	v1 "adp_backend/api/v1"
	"adp_backend/common"
	"adp_backend/config"
	"adp_backend/server"

	"github.com/dgrijalva/jwt-go"
	logger "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

const (
	_abortIndex  int8 = math.MaxInt8 / 2
	_headerAuthz      = "authorization"
	_bearer           = "Bearer"
)

var (
	_whitelist = []string{
		"/adm.ADM/Login",
		"/adm.ADM/CheckMfa",
		/*		"/adm.ADM/Logout",
				"/adm.ADM/GetLicence",
				"/adm.ADM/UpdateLicence",*/
	}
)

// ADM grpc service struct
type ADMServiceV1 struct {
	env *config.Env
}

// GrpcService is the grpc server and its configurations.
type GrpcService struct {
	env      *config.Env
	server   *grpc.Server
	handlers []grpc.UnaryServerInterceptor
}

// New news a GrpcService using customized configurations.
func New(env *config.Env, opt ...grpc.ServerOption) *GrpcService {
	keepAlive := grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle:     time.Second * 60,
		MaxConnectionAge:      time.Hour * 2,
		MaxConnectionAgeGrace: time.Second * 20,
		Time:                  time.Second * 300,
		Timeout:               time.Second * 20,
	})

	s := new(GrpcService)
	s.env = env

	opt = append(opt, keepAlive, grpc.UnaryInterceptor(s.interceptor))
	opt = append(opt, grpc.MaxRecvMsgSize(1024*1024*32))
	opt = append(opt, grpc.MaxSendMsgSize(1024*1024*32))

	s.server = grpc.NewServer(opt...)
	//s.Use(s.recovery(), s.handle(), s.logging(), s.validate())

	v1.RegisterADMServer(s.server, &ADMServiceV1{env})

	return s
}

// Start starts the grpc server.
func (s *GrpcService) Start(address string) error {
	logger.Infof("starting grpc service at: %s", address)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	go SyncComputerList(s.env)

	reflection.Register(s.server)
	return s.server.Serve(listener)
}

// Stop stops the grpc server.
func (s *GrpcService) Stop() {
	s.server.Stop()
}

// Execution is done in left-to-right order, including passing of context.
// For example ChainUnaryServer(one, two, three) will execute one before two before three, and three
// will see context changes of one and two.
func (s *GrpcService) interceptor(ctx context.Context, req interface{},
	args *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var (
		i     int
		chain grpc.UnaryHandler
	)

	n := len(s.handlers)
	if n == 0 {
		return handler(ctx, req)
	}

	chain = func(ic context.Context, ir interface{}) (interface{}, error) {
		if i == n-1 {
			return handler(ic, ir)
		}
		i++
		return s.handlers[i](ic, ir, args, chain)
	}

	return s.handlers[0](ctx, req, args, chain)
}

// Use attachs a global inteceptor to the server.
// For example, this is the right place for a rate limiter or error management inteceptor
func (s *GrpcService) Use(handlers ...grpc.UnaryServerInterceptor) *GrpcService {
	finalSize := len(s.handlers) + len(handlers)
	if finalSize >= int(_abortIndex) {
		panic("grep service: server use too many handlers")
	}
	mergedHandlers := make([]grpc.UnaryServerInterceptor, finalSize)
	copy(mergedHandlers, s.handlers)
	copy(mergedHandlers[len(s.handlers):], handlers)
	s.handlers = mergedHandlers

	return s
}

// recovery is a server interceptor that recovers from any panics.
func (s *GrpcService) recovery() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, args *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if rerr := recover(); rerr != nil {
				const size = 64 << 10
				buf := make([]byte, size)
				_ = runtime.Stack(buf, false)
				logger.Errorf("grpc server panic: %v\n%v\n%s\n", req, rerr, buf)
				err = status.Errorf(codes.Unknown, fmt.Sprintf("%v", rerr))
			}
		}()
		resp, err = handler(ctx, req)
		return
	}
}

//// handle return a new unary server interceptor for Tracing\LinkTimeout\AuthToken
func (s *GrpcService) handle() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, args *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		// check if grpc FullMethod in the whitelist
		for _, path := range _whitelist {
			if path == args.FullMethod {
				return handler(ctx, req)
			}
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.InvalidArgument, "空数据")
		}

		// 获取 header 的`authorization: Bearer token`字段的值
		var token string
		if val, ok := md[_headerAuthz]; ok {
			splits := strings.SplitN(val[0], " ", 2)
			if len(splits) < 2 || splits[0] != _bearer {
				return nil, status.Errorf(codes.Unauthenticated, "授权失败")
			}
			token = splits[1]
		}

		// 解析jwt-token进行认证
		u, err := parseToken(token, common.JWT_SECRET)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "登录过期，请重新登录")
		}

		// 检测单用户登录
		value := UserLoginCountInfo.Get(u.User)
		if value.LastLoginExpireTime > u.Expired {
			return nil, status.Errorf(codes.Unauthenticated, "已有其他用户登录，请重新登录")
		}

		// 接口鉴权
		ok, err = authentication(u, args.FullMethod)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "授权失败:%v", err)
		}
		if !ok {
			logger.Debugf("user(%s) has no permission: %s", u.User, args.FullMethod)
			return nil, status.Errorf(codes.PermissionDenied, "没有访问权限")
		}

		md.Append("token", token)
		md.Append("user", u.User)
		md.Append("role", u.Role)
		md.Append("priv", strconv.Itoa(u.Priv))
		newCtx := metadata.NewIncomingContext(ctx, md)

		return handler(newCtx, req)
	}
}

// grpc logging
func (s *GrpcService) logging() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, args *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now()
		var addr, remoteIP string
		if peerInfo, ok := peer.FromContext(ctx); ok {
			if tcpAddr, ok := peerInfo.Addr.(*net.TCPAddr); ok {
				addr = tcpAddr.IP.String()
			} else {
				addr = peerInfo.Addr.String()
			}
		}

		md, _ := metadata.FromIncomingContext(ctx) // ignore the `error` return value
		rips := md.Get("x-real-ip")
		if len(rips) > 0 {
			remoteIP = rips[0]
			addr = remoteIP
		} else {
			remoteIP = addr
		}

		var quota float64
		if deadline, ok := ctx.Deadline(); ok {
			quota = time.Until(deadline).Seconds()
		}
		// call server handler
		resp, err := handler(ctx, req)

		duration := time.Since(startTime)
		argsStr := req.(fmt.Stringer).String()
		fullMethod := args.FullMethod
		logFields := logger.Fields{
			"ip":            addr,
			"path":          fullMethod,
			"ts":            duration.Seconds(),
			"timeout_quota": quota,
			"args":          argsStr,
		}
		// add audit log
		eventResult := "成功"
		if err != nil {
			logFields["error"] = err.Error()
			logFields["stack"] = fmt.Sprintf("%+v", err)
			eventResult = "失败"
		}
		username := ""
		if len(md["user"]) > 0 {
			username = md["user"][0]
		} else if fullMethod == "/adm.ADM/Login" || fullMethod == "/adm.ADM/Logout" {
			username = getRegUser(argsStr)
		}

		if event, ok := v1.URLEventMap[args.FullMethod]; ok {
			eventArgs := eventMasking(fullMethod, argsStr)
			_ = server.AddAuditLog(s.env, username, remoteIP, event, eventArgs, eventResult)
		}
		logger.WithFields(logFields).Debugf("grpc request")
		return resp, err
	}
}

// signal handler
func (s *GrpcService) SignalHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	ch := <-c

	logger.Infof("apiserver get %s signal", ch.String())
	switch ch {
	case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
		logger.Info("apiserver exit")
		s.Stop()
		s.env.SaveProfile()
		time.Sleep(time.Second)
		return
	case syscall.SIGHUP:
		// TODO reload
	default:
		return
	}
}

// GetUser is ADMServiceV1's internal interface
func (h *ADMServiceV1) GetUser(ctx context.Context) string {
	md, _ := metadata.FromIncomingContext(ctx) // ignore the `error` return value
	return md["user"][0]
}

// GetUser is ADMServiceV1's internal interface
func (h *ADMServiceV1) IsSuper(ctx context.Context) bool {
	md, _ := metadata.FromIncomingContext(ctx) // ignore the `error` return valu
	switch md["priv"][0] {
	case strconv.Itoa(common.PrivSuper):
		return true
	default:
		return false
	}
}

// 处理err可参考: https://godoc.org/github.com/dgrijalva/jwt-go#ex-Parse--ErrorChecking
// 或jwt.MapClaims的Valid()
type UserClaim struct {
	User    string `json:"user"`
	Role    string `json:"role"`
	Priv    int    `json:"priv"`
	Expired int64  `json:"exp"`
}

func (c UserClaim) Valid() error {
	vErr := new(jwt.ValidationError)
	now := time.Now().Unix()
	if c.Expired == 0 {
		vErr.Inner = fmt.Errorf("exp is required")
		vErr.Errors |= jwt.ValidationErrorClaimsInvalid
	}
	if c.Expired < now {
		delta := time.Unix(now, 0).Sub(time.Unix(c.Expired, 0))
		vErr.Inner = fmt.Errorf("token is expired by %v", delta)
		vErr.Errors |= jwt.ValidationErrorExpired
	}

	if c.User == "" {
		vErr.Inner = fmt.Errorf("user is required")
		vErr.Errors |= jwt.ValidationErrorClaimsInvalid
	}

	if c.Role == "" {
		vErr.Inner = fmt.Errorf("role is required")
		vErr.Errors |= jwt.ValidationErrorClaimsInvalid
	}

	if c.Priv == 0 {
		vErr.Inner = fmt.Errorf("priv is required")
		vErr.Errors |= jwt.ValidationErrorClaimsInvalid
	}

	if vErr.Errors == 0 {
		return nil
	}

	return vErr
}

// 解析token获取user消息
func parseToken(tokenStr, authSecret string) (*UserClaim, error) {
	fn := func(token *jwt.Token) (interface{}, error) {
		return []byte(authSecret), nil
	}

	token, err := jwt.ParseWithClaims(tokenStr, &UserClaim{}, fn)
	if err != nil {
		return nil, err
	}

	claim, ok := token.Claims.(*UserClaim)
	if !ok {
		return nil, errors.New("cannot convert claim to BasicClaim")
	}

	return claim, nil
}

// 鉴权
func authentication(u *UserClaim, fullMethod string) (bool, error) {
	if u.Priv == common.PrivSuper {
		return true, nil
	}
	return v1.CheckUserAccess(u.Role, fullMethod), nil
}

//脱敏
func eventMasking(fullMethod string, data string) string {
	for url, mkList := range v1.URLEventMaskingMap {
		if url == fullMethod {
			for _, m := range mkList {
				reg := regexp.MustCompile(m + ":(.*)\"")
				data = reg.ReplaceAllString(data, "\""+m+"\""+`:"*""`)
			}
		}
	}
	return data
}

func getRegUser(data string) string {
	reg := regexp.MustCompile(`username:"(.*?)"`)
	return reg.FindStringSubmatch(data)[1]
}

type validator interface {
	Validate() error
}

// proto参数校验
func (s *GrpcService) validate() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if v, ok := req.(validator); ok {
			if err := v.Validate(); err != nil {
				logger.Infof("middleware validate parameter err(path:%s):%v", info.FullMethod, err)
				return nil, status.Errorf(codes.InvalidArgument, "参数错误")
			}
		}
		return handler(ctx, req)
	}
}
