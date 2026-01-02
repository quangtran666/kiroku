# Hướng dẫn sử dụng Kiroku

Chào mừng bạn đến với **Kiroku**, ứng dụng ghi chú nhanh trên terminal được viết bằng Go. Tài liệu này sẽ hướng dẫn bạn cách cài đặt, cấu hình và sử dụng các tính năng của Kiroku một cách hiệu quả nhất.

## 1. Cài đặt

### Yêu cầu

- Đã cài đặt Go (phiên bản 1.21 trở lên).
- Có `make` (nếu muốn build từ source).

### Cài đặt từ source

```bash
git clone https://github.com/tranducquang/kiroku.git
cd kiroku
make build   # Build ứng dụng
make install # Cài đặt vào GOPATH/bin
```

Hoặc cài trực tiếp bằng Go:

```bash
go install github.com/tranducquang/kiroku/cmd/kiroku@latest
```

## 2. Giao diện Terminal (TUI)

Cách chính để sử dụng Kiroku là thông qua giao diện terminal tương tác.
Chỉ cần gõ lệnh sau để mở:

```bash
kiroku
```

### Phím tắt (Keyboard Shortcuts)

Hệ thống phím tắt của Kiroku được thiết kế theo phong cách Vim để tối ưu hóa tốc độ.

#### Di chuyển & Điều hướng

| Phím      | Chức năng                        |
| --------- | -------------------------------- |
| `↑` / `k` | Di chuyển lên                    |
| `↓` / `j` | Di chuyển xuống                  |
| `←` / `h` | Quay lại / Thu gọn thư mục       |
| `→` / `l` | Chọn / Mở rộng thư mục           |
| `Tab`     | Chuyển đổi giữa các bảng (Panel) |
| `Enter`   | Chọn / Xác nhận                  |
| `Esc`     | Quay lại / Hủy bỏ                |

#### Thao tác chính

| Phím | Chức năng                    |
| ---- | ---------------------------- |
| `n`  | Tạo ghi chú mới (New Note)   |
| `t`  | Tạo Todo mới (New Todo)      |
| `N`  | Tạo thư mục mới (New Folder) |
| `e`  | Chỉnh sửa nội dung (Edit)    |
| `d`  | Xóa (Delete)                 |
| `/`  | Tìm kiếm (Search)            |

#### Quản lý & Tiện ích

| Phím           | Chức năng                                  |
| -------------- | ------------------------------------------ |
| `s`            | Đánh dấu sao / Bỏ đánh dấu (Star)          |
| `x` / `Space`  | Đánh dấu hoàn thành Todo (Toggle Done)     |
| `p`            | Thay đổi độ ưu tiên Todo (Priority)        |
| `m`            | Di chuyển ghi chú sang thư mục khác (Move) |
| `P`            | Bật/tắt chế độ xem trước (Preview)         |
| `r` / `Ctrl+r` | Tải lại dữ liệu (Refresh)                  |
| `?`            | Xem trợ giúp (Help)                        |
| `q` / `Ctrl+c` | Thoát ứng dụng                             |

## 3. Sử dụng dòng lệnh (CLI)

Ngoài giao diện TUI, bạn có thể dùng các lệnh CLI để thực hiện nhanh các tác vụ mà không cần mở giao diện đầy đủ.

### Tạo ghi chú nhanh

```bash
# Ghi chú cơ bản
kiroku add "Tiêu đề ghi chú"

# Vào thư mục cụ thể
kiroku add "Báo cáo ngày" -f work

# Dùng template có sẵn
kiroku add "Họp team" -t meeting-notes
```

### Tạo Todo nhanh

```bash
# Todo cơ bản
kiroku todo "Mua cà phê"

# Có độ ưu tiên cao
kiroku todo "Deadline dự án" -p high

# Có ngày hết hạn
kiroku todo "Nộp báo cáo" -d 2026-01-05
```

### Tìm kiếm & Liệt kê

```bash
# Liệt kê tất cả
kiroku list

# Chỉ xem việc cần làm (todos)
kiroku list --todos

# Tìm kiếm nội dung
kiroku list --search "golang"
```

## 4. Cấu hình

File cấu hình nằm tại: `~/.config/kiroku/config.yaml`
Bạn có thể trỏ đến file database, chọn trình chỉnh sửa văn bản (editor) ưa thích và tùy chỉnh giao diện.

Ví dụ file cấu hình chuẩn:

```yaml
database:
  path: ~/.local/share/kiroku/kiroku.db

editor:
  command: nvim # Hoặc 'vim', 'code', 'nano'
  args: ["-c", "set filetype=markdown"]

ui:
  theme: dark
  show_preview: true
  sidebar_width: 25

todos:
  show_completed: true
  sort_by: priority
```

## 5. Mẹo sử dụng

1. **Templates**: Tận dụng tính năng Template để tạo nhanh các mẫu ghi chú lặp lại (Daily Standup, Bug Report).
2. **Vim Mode**: Nếu bạn quen dùng Vim, việc điều hướng bằng `h/j/k/l` sẽ cực kỳ mượt mà.
3. **Todo Priorities**: Sử dụng phím `p` để xoay vòng các mức độ ưu tiên của Todo (Low -> Medium -> High).

Chúc bạn ghi chú hiệu quả với Kiroku!
