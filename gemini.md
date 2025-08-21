# Trình quản lý kết nối SSH - Tài liệu dự án

## Tổng quan dự án

**Tên dự án:** Trình quản lý kết nối SSH (SSH Connection Manager)
**Ngôn ngữ:** Go (Golang)
**Loại:** Công cụ dòng lệnh (CLI)
**Nền tảng:** Đa nền tảng (Linux, macOS, Windows)
**Phân phối:** Tệp thực thi đơn

### Mục đích chính
Một công cụ dòng lệnh để quản lý các kết nối SSH một cách hiệu quả, cho phép người dùng lưu, sắp xếp và kết nối nhanh đến các máy chủ từ xa mà không cần nhớ các lệnh và cấu hình SSH phức tạp.

## Cấu trúc dự án

```
ssh-manager/
├── cmd/
│   ├── root.go              # Cài đặt lệnh gốc
│   ├── add.go               # Thêm kết nối SSH mới
│   ├── list.go              # Liệt kê tất cả kết nối
│   ├── connect.go           # Kết nối đến máy chủ đã lưu
│   ├── remove.go            # Xóa kết nối
│   ├── edit.go              # Chỉnh sửa kết nối hiện có
│   └── export.go            # Xuất/nhập cấu hình
├── internal/
│   ├── config/
│   │   ├── config.go        # Quản lý cấu hình
│   │   └── storage.go       # Lớp lưu trữ dữ liệu
│   ├── ssh/
│   │   ├── connection.go    # Logic kết nối SSH
│   │   ├── client.go        # Wrapper cho SSH client
│   │   └── keys.go          # Quản lý khóa SSH
│   ├── models/
│   │   └── connection.go    # Mô hình dữ liệu
│   └── utils/
│       ├── crypto.go        # Mã hóa dữ liệu nhạy cảm
│       ├── validation.go    # Xác thực đầu vào
│       └── terminal.go      # Các hàm trợ giúp cho giao diện terminal
├── pkg/
│   └── version/
│       └── version.go       # Thông tin phiên bản
├── configs/
│   └── example.yaml         # Cấu hình mẫu
├── scripts/
│   ├── build.sh            # Script build cho tất cả nền tảng
│   ├── install.sh          # Script cài đặt
│   └── release.sh          # Tự động hóa việc phát hành
├── docs/
│   ├── USAGE.md            # Tài liệu hướng dẫn sử dụng
│   ├── CONFIGURATION.md    # Hướng dẫn cấu hình
│   └── DEVELOPMENT.md      # Hướng dẫn cài đặt môi trường phát triển
├── .github/
│   └── workflows/
│       └── release.yml     # GitHub Actions cho việc phát hành
├── go.mod
├── go.sum
├── main.go                 # Điểm vào của ứng dụng
├── Makefile               # Tự động hóa build
├── README.md              # Tệp README của dự án
└── LICENSE                # Giấy phép MIT
```

## Yêu cầu kỹ thuật

### Dependencies (Các thư viện phụ thuộc)
- **CLI Framework:** `github.com/spf13/cobra` - Framework CLI hiện đại
- **Configuration:** `github.com/spf13/viper` - Quản lý cấu hình
- **SSH Client:** `golang.org/x/crypto/ssh` - Triển khai SSH chính thức
- **Encryption:** `golang.org/x/crypto` - Để lưu trữ thông tin đăng nhập an toàn
- **Terminal UI:** `github.com/manifoldco/promptui` - Các câu lệnh tương tác
- **File Operations:** Thư viện chuẩn `os`, `filepath`
- **JSON/YAML:** `gopkg.in/yaml.v3` cho các tệp cấu hình

### Phiên bản Go
- **Tối thiểu:** Go 1.19+
- **Khuyến nghị:** Go 1.21+

## Đặc tả tính năng cốt lõi

### 1. Quản lý kết nối
```bash
# Thêm kết nối mới
ssh-manager add <name> --host <host> --user <user> [--port <port>] [--key <path>] [--pass <password>]

# Liệt kê tất cả kết nối
ssh-manager list [--format table|json]

# Kết nối đến máy chủ
ssh-manager connect <name>
ssh-manager <name>  # viết tắt

# Xóa kết nối
ssh-manager remove <name>

# Chỉnh sửa kết nối hiện có
ssh-manager edit <name>
```

### 2. Quản lý cấu hình
```bash
# Hiển thị cấu hình
ssh-manager config show

# Đặt giá trị mặc định toàn cục
ssh-manager config set default-user myuser
ssh-manager config set default-port 2222

# Xuất/Nhập
ssh-manager export --output backup.yaml
ssh-manager import --input backup.yaml
```

### 3. Quản lý khóa SSH
```bash
# Liệt kê các khóa SSH
ssh-manager keys list

# Thêm khóa SSH
ssh-manager keys add --name work --path ~/.ssh/work_rsa

# Tạo cặp khóa mới
ssh-manager keys generate --name newkey --type rsa --bits 4096
```

## Mô hình dữ liệu

### Mô hình Connection
```go
type Connection struct {
    Name        string            `json:"name" yaml:"name"`
    Host        string            `json:"host" yaml:"host"`
    Port        int               `json:"port" yaml:"port"`
    User        string            `json:"user" yaml:"user"`
    KeyPath     string            `json:"key_path,omitempty" yaml:"key_path,omitempty"`
    Password    string            `json:"password,omitempty" yaml:"password,omitempty"` // đã mã hóa
    Tags        []string          `json:"tags,omitempty" yaml:"tags,omitempty"`
    Description string            `json:"description,omitempty" yaml:"description,omitempty"`
    LastUsed    time.Time         `json:"last_used" yaml:"last_used"`
    CreatedAt   time.Time         `json:"created_at" yaml:"created_at"`
    Extra       map[string]string `json:"extra,omitempty" yaml:"extra,omitempty"`
}
```

### Mô hình Config
```go
type Config struct {
    DefaultUser     string                `yaml:"default_user"`
    DefaultPort     int                   `yaml:"default_port"`
    DefaultKeyPath  string                `yaml:"default_key_path"`
    ConfigPath      string                `yaml:"config_path"`
    Connections     map[string]Connection `yaml:"connections"`
    SSHKeys         map[string]SSHKey     `yaml:"ssh_keys"`
    Settings        Settings              `yaml:"settings"`
}

type Settings struct {
    EncryptPasswords bool   `yaml:"encrypt_passwords"`
    LogConnections   bool   `yaml:"log_connections"`
    LogPath         string `yaml:"log_path"`
    Editor          string `yaml:"editor"`
}
```

## Yêu cầu bảo mật

### 1. Mã hóa mật khẩu
- Sử dụng AES-256-GCM để mã hóa mật khẩu
- Lưu trữ khóa mã hóa trong system keyring (Linux: libsecret, macOS: Keychain, Windows: Credential Manager)
- Không bao giờ lưu trữ mật khẩu dưới dạng văn bản thuần túy

### 2. Bảo mật khóa SSH
- Hỗ trợ các loại khóa RSA, ECDSA, Ed25519
- Xác thực quyền của tệp khóa (600/400)
- Cảnh báo về các khóa có thể đọc được bởi mọi người

### 3. Bảo mật cấu hình
- Quyền tệp cấu hình: 600 (chỉ chủ sở hữu có quyền đọc/ghi)
- Tạo tệp tạm thời an toàn
- Xóa dữ liệu nhạy cảm khỏi bộ nhớ sau khi sử dụng

## Nguyên tắc thiết kế CLI

### 1. Tính khả dụng
- Cấu trúc lệnh trực quan
- Các câu lệnh tương tác cho thông tin còn thiếu
- Thông báo lỗi hữu ích kèm theo đề xuất
- Hỗ trợ tự động hoàn thành (tab completion)

### 2. Định dạng đầu ra
- Mặc định: Định dạng bảng dễ đọc cho người dùng
- JSON: Để viết kịch bản và tích hợp
- YAML: Để xuất cấu hình
- Đầu ra có màu (tắt bằng `--no-color`)

### 3. Cấu hình
- Vị trí tệp cấu hình (theo thứ tự ưu tiên):
  1. Cờ `--config`
  2. Biến môi trường `$SSH_MANAGER_CONFIG`
  3. `$HOME/.ssh-manager/config.yaml`
  4. `$HOME/.config/ssh-manager/config.yaml`

## Đặc tả Build và Release

### Mục tiêu Build
```makefile
# Makefile targets
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

build-all:
	@for platform in $(PLATFORMS); do \
		GOOS=$${platform%/*} GOARCH=$${platform#*/} \
		go build -ldflags "-X main.version=$(VERSION)" \
		-o bin/ssh-manager-$${platform%/*}-$${platform#*/} .; \
	done
```

### Chiến lược phát hành
1. **GitHub Releases:** Tự động hóa qua GitHub Actions
2. **Trình quản lý gói:**
   - Homebrew formula cho macOS
   - Gói AUR cho Arch Linux
   - Gói .deb cho Ubuntu/Debian
3. **Tải xuống trực tiếp:** Các tệp nhị phân được build sẵn cho tất cả các nền tảng

## Hướng dẫn phát triển

### Phong cách code
- Tuân thủ các quy ước của Go (`gofmt`, `golint`)
- Sử dụng tên biến có ý nghĩa
- Viết test toàn diện (mục tiêu: độ bao phủ >80%)
- Viết tài liệu cho tất cả các hàm được export

### Xử lý lỗi
- Sử dụng wrapped errors với ngữ cảnh
- Cung cấp thông báo lỗi có thể hành động
- Ghi log lỗi một cách thích hợp
- Suy giảm chức năng một cách duyên dáng khi có thể

### Chiến lược kiểm thử
- Unit test cho tất cả logic cốt lõi
- Integration test cho các kết nối SSH
- CLI test sử dụng golden files
- Kiểm thử đa nền tảng qua CI

## Yêu cầu trải nghiệm người dùng

### 1. Cài đặt lần đầu
- Trình hướng dẫn cài đặt tương tác: `ssh-manager init`
- Nhập từ cấu hình SSH hiện có
- Tạo cấu hình mẫu

### 2. Khả năng khám phá
- Hệ thống trợ giúp tích hợp
- Gợi ý lệnh cho các lỗi gõ sai
- Ví dụ trong văn bản trợ giúp

### 3. Tính năng tăng năng suất
- Tìm kiếm mờ cho tên kết nối
- Danh sách các kết nối gần đây
- Thao tác hàng loạt
- Tự động hoàn thành cho shell

## Yêu cầu về hiệu năng

### 1. Thời gian khởi động
- Khởi động nguội: < 100ms
- Kích thước tệp nhị phân: < 25MB
- Mức sử dụng bộ nhớ: < 10MB khi khởi động

### 2. Tốc độ kết nối
- Thiết lập kết nối: < 2 giây đối với mạng cục bộ
- Tải cấu hình: < 50ms
- Thao tác liệt kê: < 100ms cho hơn 1000 kết nối

## Yêu cầu về tương thích

### Hệ điều hành
- Linux: Ubuntu 18.04+, RHEL 7+, Arch Linux
- macOS: 10.15+ (Catalina và mới hơn)
- Windows: Windows 10+

### Tương thích SSH
- Giao thức SSH phiên bản 2.0
- Hỗ trợ các thuật toán trao đổi khóa phổ biến
- Tương thích với các cấu hình máy chủ OpenSSH

## Yêu cầu về tài liệu

### Tài liệu người dùng
- `README.md`: Bắt đầu nhanh và tổng quan
- `docs/USAGE.md`: Hướng dẫn sử dụng toàn diện
- `docs/CONFIGURATION.md`: Tham chiếu cấu hình
- Tạo trang man

### Tài liệu cho nhà phát triển
- `docs/DEVELOPMENT.md`: Hướng dẫn cài đặt và đóng góp
- Bình luận code cho tất cả các API công khai
- Ghi lại các quyết định kiến trúc (ADRs)

## Các chỉ số thành công

### Chức năng
- [ ] Tất cả các tính năng cốt lõi được triển khai và kiểm thử
- [ ] Khả năng tương thích đa nền tảng được xác minh
- [ ] Các yêu cầu bảo mật được đáp ứng
- [ ] Các mục tiêu về hiệu năng đã đạt được

### Chất lượng
- [ ] Độ bao phủ test > 80%
- [ ] Không có lỗ hổng bảo mật nghiêm trọng
- [ ] Tài liệu đầy đủ và chính xác
- [ ] Phản hồi của người dùng được tích hợp

### Phân phối
- [ ] Phát hành tự động hoạt động
- [ ] Tích hợp trình quản lý gói
- [ ] Các kịch bản cài đặt đã được kiểm thử
- [ ] Lộ trình nâng cấp được ghi lại

## Các giai đoạn triển khai

### Giai đoạn 1: Cơ sở hạ tầng cốt lõi (Tuần 1)
- Thiết lập cấu trúc dự án
- Framework CLI cơ bản
- Quản lý cấu hình
- Mô hình dữ liệu

### Giai đoạn 2: Quản lý kết nối (Tuần 2)
- Thêm/xóa/liệt kê kết nối
- Chức năng kết nối SSH cơ bản
- I/O tệp cấu hình
- Xác thực đầu vào

### Giai đoạn 3: Các tính năng nâng cao (Tuần 3)
- Quản lý khóa SSH
- Mã hóa mật khẩu
- Các câu lệnh tương tác
- Chức năng xuất/nhập

### Giai đoạn 4: Hoàn thiện và phát hành (Tuần 4)
- Kiểm thử toàn diện
- Hoàn thành tài liệu
- Tự động hóa build
- Chuẩn bị phát hành

## Hướng dẫn cho trợ lý AI

Khi làm việc với dự án này:

1. **Tuân thủ cấu trúc**: Triển khai các tệp theo cấu trúc dự án đã chỉ định
2. **Triển khai tăng dần**: Bắt đầu với chức năng cốt lõi, sau đó thêm các tính năng
3. **Kiểm thử kỹ lưỡng**: Viết test cho từng thành phần khi bạn xây dựng
4. **Viết tài liệu song song**: Cập nhật tài liệu khi thêm tính năng
5. **Bảo mật là trên hết**: Luôn xem xét các tác động về bảo mật
6. **Trải nghiệm người dùng**: Suy nghĩ từ góc độ của người dùng để thiết kế CLI
7. **Đa nền tảng**: Kiểm thử và xem xét sự khác biệt giữa các hệ điều hành
8. **Hiệu năng**: Giữ cho công cụ nhanh và nhẹ

Khi được yêu cầu triển khai các thành phần cụ thể, hãy tham khảo tài liệu này để biết ngữ cảnh, yêu cầu và các quyết định kiến trúc.
