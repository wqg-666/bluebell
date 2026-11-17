package logger

import (
	settings "bluebell/setting"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LogConfig 这个结构体用来存放日志配置
type LogConfig struct {
	Filename   string `json:"filename" yaml:"filename"`
	MaxSize    int    `json:"maxsize" yaml:"maxsize"`
	MaxBackups int    `json:"maxbackups" yaml:"maxbackups"`
	MaxAge     int    `json:"maxage" yaml:"maxage"`
	Level      string `json:"level" yaml:"level"`
}

// Init 初始化 Zap 日志系统
func Init(cfg *settings.LogConfig, mode string) (err error) {
	//日志配置结构体 的指针，里面保存日志文件路径、大小、等级等。
	//mode：运行模式，例如 dev 表示开发模式、

	//调用 getLogWriter()，创建一个日志写入器。它负责决定日志写入哪里，以及日志文件如何切割
	writeSyncer := getLogWriter(
		cfg.Filename,
		cfg.MaxSize,
		cfg.MaxBackups,
		cfg.MaxAge,
	)
	//调用 getEncoder()，决定日志以什么格式输出。这里采用 JSON 格式。
	encoder := getEncoder()

	var l = new(zapcore.Level)               //创建一个 zapcore.Level 类型的指针。
	err = l.UnmarshalText([]byte(cfg.Level)) //将配置文件中的字符串转换为 Zap 能识别的日志等级。
	if err != nil {
		return
	}
	//lumberjack：管理日志文件。
	//WriteSyncer：负责日志写入。
	//Encoder：负责日志格式。
	//Core：把它们组合起来。zapcore.Core连接日志等级、编码器和写入器

	//Core 是 Zap 日志系统的核心组件。
	//它主要负责：
	//- 判断日志等级是否需要输出。
	//- 使用 Encoder 编码日志。
	//- 使用 WriteSyncer 写入日志。
	var core zapcore.Core
	if mode == "dev" {
		//决定控制台中的日志以什么格式显示。
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		core = zapcore.NewTee( //zapcore.NewTee(...)就是把多个 Core 合并起来，让一条日志同时交给多个 Core 处理。
			//开发模式下，一条日志既会写入文件，也会输出到控制台。
			//日志文件
			zapcore.NewCore(encoder, writeSyncer, l),
			//控制台
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel),
			//consoleEncoder：控制台日志编码格式。
			//os.Stdout：标准输出，也就是终端。
			//zapcore.Lock()：给写入器加锁，避免并发写入时出现问题。
			//zapcore.DebugLevel：允许 Debug 及以上等级的日志输出。
		)
	} else {
		//日志文件
		core = zapcore.NewCore(encoder, writeSyncer, l)
	}
	//根据刚才创建的 Core，创建一个 Logger。这里创建了一个 Logger 对象，并赋值给变量 lg。
	//zap.AddCaller()表示记录日志调用的位置，这样出现错误时，就能知道日志是从哪个文件的哪一行代码输出的
	lg := zap.New(core, zap.AddCaller())
	//把自己创建的logger替换zap库中全局的logger
	zap.ReplaceGlobals(lg)
	//对比
	//替换，当成全局logger使用
	//zap.L().Info("服务器启动成功")，这里的 L() 是 Zap 提供的一个函数，用于获取当前的全局 Logger。
	//如果不替换
	//直接使用 lg
	// lg.Info("服务器启动成功")
	//func Login() {
	//    zap.L().Info("用户登录")
	//}
	//func Login(lg *zap.Logger) {
	//    lg.Info("用户登录")
	//}
	return
}

// 设置日志输出格式
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()             //获取 Zap 默认的生产环境编码配置。
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder         //设置时间格式
	encoderConfig.TimeKey = "time"                                //修改时间字段名称
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder       //设置日志等级格式，将日志等级转换成大写。
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder //设置耗时格式，将耗时转换成秒。
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder       //设置调用位置格式，只保留简短的文件名和行号
	return zapcore.NewJSONEncoder(encoderConfig)                  //创建 JSON 编码器。
}

// 设置日志写入位置和切割规则
func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	//第三方库：lumberjack它的主要作用是：管理日志文件，支持日志切割和旧日志清理。
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,   //单个文件最大 10MB。
		MaxBackups: maxBackup, //最多保留 5 个旧文件。
		MaxAge:     maxAge,    //旧日志最多保留的天。
	}
	return zapcore.AddSync(lumberJackLogger) //把lumberjack包装成 Zap 需要的 WriteSyncer 类型。
}

// GinLogger 记录 Gin 请求日志。前面都是初始化日志系统，这里开始与 Gin 框架结合。
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()             //记录请求开始的时间。
		path := c.Request.URL.Path      //获取请求路径
		query := c.Request.URL.RawQuery //获取 URL 查询参数
		//例如：http://localhost:8080/user/list?id=10&name=张三
		//那么：query = "id=10&name=张三"
		c.Next()                  //继续执行后面的 Gin 中间件以及最终的路由处理函数
		cost := time.Since(start) //计算请求耗时
		//记录请求信息，生成的日志格式
		//这里的 L() 是 Zap 提供的一个函数，用于获取当前的全局 Logger。
		zap.L().Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()), //记录客户端 IP
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

// GinRecovery 捕获 Gin 请求中的 panic
// 1. 遇到 panic。 2. 开始执行 defer。
// 3. recover() 捕获 panic。 4. 程序不会因为这次 panic 直接崩溃。
// 这里的 GinRecovery 就是利用这个机制保护请求处理过程。
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			//recover()：捕获当前 Goroutine 中正在发生的 panic，阻止其继续向上传播。
			if err := recover(); err != nil {
				//判断是否是连接断开异常
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				//获取 HTTP 请求信息，这个函数可以把 HTTP 请求转换成字节数据，方便写入日志。第二个参数 false 表示不把请求体一起导出。
				//User-Agent: Mozilla/5.0
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				//如果是连接断开
				if brokenPipe {
					//zap.L()：获取 Zap 的全局 Logger。
					//记录错误日志
					zap.L().Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					//将错误添加到 Gin Context
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					//终止后续中间件，阻止后续 Handler 继续执行。
					c.Abort()
					return
				}
				//判断是否需要记录堆栈信息，这里的 stack 是 GinRecovery 的参数。
				//如果： GinRecovery(true)就记录堆栈信息。
				if stack {
					zap.L().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					zap.L().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError) //返回 500 状态码
			}
		}()
		c.Next()
	}
}
