package core

import (
	"encoding/json"
	"strings"
	"sync"
	"app/cau_hinh"
)

const (
	DongBatDau_TinNhan = 11

	CotTN_MaTinNhan    = 0  // A
	CotTN_LoaiTinNhan  = 1  // B
	CotTN_NguoiGuiID   = 2  // C
	CotTN_NguoiNhanID  = 3  // D
	CotTN_TieuDe       = 4  // E
	CotTN_NoiDung      = 5  // F
	CotTN_DinhKemJson  = 6  // G
	CotTN_ThamChieuID  = 7  // H
	CotTN_ReplyChoID   = 8  // I
	CotTN_NgayTao      = 9  // J
	CotTN_NguoiDocJson = 10 // K
	CotTN_TrangThaiXoa = 11 // L
	// Mở rộng thêm 2 cột theo thiết kế hiển thị
	CotTN_TenNguoiGui    = 12 // M
	CotTN_ChucVuNguoiGui = 13 // N
)

type FileDinhKem struct {
	TenFile string `json:"name"`
	URL     string `json:"url"`
	Loai    string `json:"type"` 
}

type TinNhan struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`

	MaTinNhan      string        `json:"ma_tin_nhan"`
	LoaiTinNhan    string        `json:"loai_tin_nhan"`
	NguoiGuiID     string        `json:"nguoi_gui_id"`
	NguoiNhanID    string        `json:"nguoi_nhan_id"`
	TieuDe         string        `json:"tieu_de"`
	NoiDung        string        `json:"noi_dung"`
	DinhKem        []FileDinhKem `json:"dinh_kem"`
	ThamChieuID    string        `json:"tham_chieu_id"`
	ReplyChoID     string        `json:"reply_cho_id"`
	NgayTao        string        `json:"ngay_tao"`
	NguoiDoc       []string      `json:"nguoi_doc"`
	TrangThaiXoa   []string      `json:"trang_thai_xoa"`
	
	// Thông tin phụ trợ hiển thị nhanh
	TenNguoiGui    string `json:"ten_nguoi_gui"`
	ChucVuNguoiGui string `json:"chuc_vu_nguoi_gui"`
}

var (
	CacheTinNhan = make(map[string][]*TinNhan)
	mtxTinNhan   sync.RWMutex
)

// Helper check chuỗi trong mảng
func ContainsString(slice []string, val string) bool {
	for _, item := range slice {
		if item == val { return true }
	}
	return false
}

func NapTinNhan(shopID string) {
	if shopID == "" { shopID = cau_hinh.BienCauHinh.IdFileSheet }
	raw, err := LoadSheetData(shopID, "TIN_NHAN")
	if err != nil { return }

	list := []*TinNhan{}
	for i, r := range raw {
		if i < DongBatDau_TinNhan-1 { continue }
		maTN := LayString(r, CotTN_MaTinNhan)
		if maTN == "" { continue }

		tn := &TinNhan{
			SpreadsheetID:  shopID,
			DongTrongSheet: i + 1,
			MaTinNhan:      maTN,
			LoaiTinNhan:    LayString(r, CotTN_LoaiTinNhan),
			NguoiGuiID:     LayString(r, CotTN_NguoiGuiID),
			NguoiNhanID:    LayString(r, CotTN_NguoiNhanID),
			TieuDe:         LayString(r, CotTN_TieuDe),
			NoiDung:        LayString(r, CotTN_NoiDung),
			ThamChieuID:    LayString(r, CotTN_ThamChieuID),
			ReplyChoID:     LayString(r, CotTN_ReplyChoID),
			NgayTao:        LayString(r, CotTN_NgayTao),
			TenNguoiGui:    LayString(r, CotTN_TenNguoiGui),
			ChucVuNguoiGui: LayString(r, CotTN_ChucVuNguoiGui),
		}

		_ = json.Unmarshal([]byte(LayString(r, CotTN_DinhKemJson)), &tn.DinhKem)
		_ = json.Unmarshal([]byte(LayString(r, CotTN_NguoiDocJson)), &tn.NguoiDoc)
		_ = json.Unmarshal([]byte(LayString(r, CotTN_TrangThaiXoa)), &tn.TrangThaiXoa)

		if tn.DinhKem == nil { tn.DinhKem = make([]FileDinhKem, 0) }
		if tn.NguoiDoc == nil { tn.NguoiDoc = make([]string, 0) }
		if tn.TrangThaiXoa == nil { tn.TrangThaiXoa = make([]string, 0) }

		list = append(list, tn)
	}

	mtxTinNhan.Lock()
	CacheTinNhan[shopID] = list
	mtxTinNhan.Unlock()
}

// Lấy hộp thư của 1 người (Đã lọc theo ID và Vai trò)
func LayHopThuNguoiDung(shopID string, maKH string, vaiTro string) []*TinNhan {
	mtxTinNhan.RLock()
	defer mtxTinNhan.RUnlock()

	var inbox []*TinNhan
	for _, m := range CacheTinNhan[shopID] {
		// Bỏ qua tin nhắn bị user này xóa mềm
		if ContainsString(m.TrangThaiXoa, maKH) { continue }

		isReceiver := false
		
		// 1. Nếu là tin nhắn CHAT hoặc AUTO (đích danh 1 ID)
		if m.LoaiTinNhan == "CHAT" || m.LoaiTinNhan == "AUTO" {
			if m.NguoiNhanID == maKH || m.NguoiGuiID == maKH {
				isReceiver = true
			}
		} else if m.LoaiTinNhan == "SYSTEM" || m.LoaiTinNhan == "ALL" {
			// 2. Nếu là thông báo tập thể (Mảng JSON các ID hoặc Role)
			if strings.Contains(m.NguoiNhanID, "["+maKH+"]") || strings.Contains(m.NguoiNhanID, "\""+maKH+"\"") {
				isReceiver = true
			} else if m.NguoiNhanID == "ALL" || strings.Contains(m.NguoiNhanID, vaiTro) {
				isReceiver = true
			} else if m.NguoiGuiID == maKH {
				isReceiver = true // Người gửi cũng được xem lại thông báo của mình
			}
		}

		if isReceiver { inbox = append(inbox, m) }
	}
	return inbox
}

// =========================================================
// [SỬ DỤNG HÀNG ĐỢI UPDATE]: Chỉ cập nhật 1 ô (Cột NguoiDoc)
// =========================================================
func DanhDauDocTinNhan(shopID string, maKH string, msgID string) {
	mtxTinNhan.Lock()
	defer mtxTinNhan.Unlock()
	
	for _, m := range CacheTinNhan[shopID] {
		if m.MaTinNhan == msgID {
			if !ContainsString(m.NguoiDoc, maKH) {
				m.NguoiDoc = append(m.NguoiDoc, maKH)
				// Hàm cũ: Ghi đè vào đúng tọa độ dòng/cột
				ThemVaoHangCho(shopID, "TIN_NHAN", m.DongTrongSheet, CotTN_NguoiDocJson, ToJSON(m.NguoiDoc))
			}
			break
		}
	}
}

// =========================================================
// [SỬ DỤNG HÀNG ĐỢI APPEND]: Tạo tin nhắn mới hoàn toàn
// =========================================================
func ThemMoiTinNhan(shopID string, msg *TinNhan) {
	mtxTinNhan.Lock()
	list := CacheTinNhan[shopID]
	
	// 1. Dự đoán RowIndex trên RAM để sau này Update (đọc/xoá) có cái mà dùng
	maxRow := DongBatDau_TinNhan - 1
	for _, m := range list {
		if m.DongTrongSheet > maxRow { maxRow = m.DongTrongSheet }
	}
	msg.DongTrongSheet = maxRow + 1
	
	CacheTinNhan[shopID] = append(list, msg)
	mtxTinNhan.Unlock()
	
	// Khởi tạo mảng rỗng để chống nil khi parse JSON
	if msg.DinhKem == nil { msg.DinhKem = make([]FileDinhKem, 0) }
	if msg.NguoiDoc == nil { msg.NguoiDoc = make([]string, 0) }
	if msg.TrangThaiXoa == nil { msg.TrangThaiXoa = make([]string, 0) }

	// 2. Gói toàn bộ dữ liệu thành 1 dòng (Row Slice)
	dongMoi := []interface{}{
		msg.MaTinNhan,        // A (0)
		msg.LoaiTinNhan,      // B (1)
		msg.NguoiGuiID,       // C (2)
		msg.NguoiNhanID,      // D (3)
		msg.TieuDe,           // E (4)
		msg.NoiDung,          // F (5)
		ToJSON(msg.DinhKem),  // G (6)
		msg.ThamChieuID,      // H (7)
		msg.ReplyChoID,       // I (8)
		msg.NgayTao,          // J (9)
		ToJSON(msg.NguoiDoc), // K (10)
		ToJSON(msg.TrangThaiXoa), // L (11)
		msg.TenNguoiGui,      // M (12)
		msg.ChucVuNguoiGui,   // N (13)
	}

	// 3. Ném vào Hàng đợi Append (Để API Google tự tìm dòng trống cuối cùng và điền vào)
	ThemDongVaoHangCho(shopID, "TIN_NHAN", dongMoi)
}
