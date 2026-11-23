Golang 工具包（版本 >= 1.22）

### 并发控制 `sync.go`
- `RunAll` 总执行 m 个任务，并发 n 个
- `NewPool` 对象池 (sync.Pool的泛型封装)
- `NewOnce` 单例懒加载容器 （sync.Once的泛型封装）


### 通用 `common.go`
- `MustJSON` 任意值序列化为 JSON 字符串 (不可json化的）
- `MustRaw` 任意值序列化为 json.RawMessage，失败返回 nil
- `DivTrunc` 安全除法，保留x位小数（非法值处理为返回0）
- `Trunc` 截断浮点数到指定小数位（非四舍五入）
- `Sha256` 加密
- `HmacSha256` 加密
- `GetCaller` 获取调用方的文件、行号与函数名

### 检查 `check.go`
- `IsASCII` 字符串是否全部由 ASCII 字符组成。
- `IsNumber` 浮点数是否有效（浮点数可能是NaN和无穷）
- `IsPhoneNumber` 是否为国内11位手机号
- `IsLocalIP` 是否为本地/内网 IP

### 调试 
- `IsDebug` 编译期常量，默认 false，加 -tags=debug 编译时为 true，用于死码消除

### 文件操作 `file.go`
- `IsFileExist` 判断文件/文件夹是否存在（不校验权限）
- `ReadFileByLine` 按行读取文件并回调处理（单行最大 16MB）
- `ReadCSV` 逐行读取 CSV 并回调处理

### HTTP 请求
- `HttpGet` GET 请求
- `HttpPostJson` POST JSON 请求
- `HttpPostForm` POST FORM 请求

### 日志
- `NewLogger` 创建新的日志输出
- `SetLogger` 设置默认日志输出位置
- `LogDebug` 默认日志打印 Debug 级别日志
- `LogInfo` 默认日志打印 Info 级别日志
- `LogWarn` 默认日志打印 Warn 级别日志
- `LogError` 默认日志打印 Error 级别日志，并附带调用位置

### 随机
- `RandomIntN` 生成 `[0,max)` 随机数
- `RandomString` 随机字符串

