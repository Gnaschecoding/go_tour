package setting

import "time"

type ServerSettings struct {
	RunMode      string
	HttpPort     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type AppSettings struct {
	DefaultPageSize       int
	MaxPageSize           int
	DefaultContextTimeout time.Duration
	LogSavePath           string
	LogFileName           string
	LogFileExt            string

	UploadSavePath       string
	UploadServerUrl      string
	UploadImageMaxSize   int
	UploadImageAllowExts []string
}

type DatabaseSettings struct {
	DBType       string
	UserName     string
	Password     string
	Host         string
	DBName       string
	TablePrefix  string
	Charset      string
	ParseTime    bool
	MaxIdleConns int
	MaxOpenConns int
}

type JWTSettings struct {
	Secret string
	Issuer string
	Expire time.Duration
}

type EmailSettings struct {
	Host     string
	Port     int
	Username string
	Password string
	IsSSL    bool
	From     string
	To       []string
}

type RedisSettings struct {
	Network      string        // 网络类型，通常为 "tcp"
	Addr         string        // Redis 服务器地址，格式为 "host:port"
	Username     string        // Redis 认证的用户名
	Password     string        // Redis 认证的密码
	DialTimeout  time.Duration // 建立连接的超时时间
	ReadTimeout  time.Duration // 读取操作的超时时间
	WriteTimeout time.Duration // 写入操作的超时时间
	PoolSize     int           // 连接池的最大连接数
	MinIdleConns int           // 连接池中的最小空闲连接数
	MaxIdleConns int           // 连接池中的最大空闲连接数
	MaxRetries   int           // 操作失败后的最大重试次数
	DB           int           //表示连接哪个数据库
}

type LimiterSettings struct {
	FillInterval time.Duration
	Quantum      int64
	Capacity     int64
	Expiration   time.Duration //自动过期时间
}

var sections = make(map[string]interface{})

func (s *Setting) ReadSection(k string, v interface{}) error {
	err := s.vp.UnmarshalKey(k, v) //是把yaml中的关键字中的内容，（读出来的应该是json类型） 反序列化到 对应的结构体对象中，这里是用interface来接收
	if err != nil {
		return err
	}

	if _, ok := sections[k]; !ok {
		sections[k] = v
	}

	return nil
}

func (s *Setting) ReloadAllSection() error {
	for k, v := range sections {
		err := s.ReadSection(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}
