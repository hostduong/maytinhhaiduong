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
)

type FileDinhKem struct {
	TenFile string `json:"name"`
	URL     string `json:"url"`
	Loai    string `json:"type"` 
}

type TinNhan struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`

	MaTinNhan    string        `json:"ma_tin_nhan"`
	LoaiTinNhan  string        `json:"loai_tin_nhan"` // ALL, AUTO, CHAT
	NguoiGuiID   string        `json:"nguoi_gui_id"`
	NguoiNhanID  string        `json:"nguoi_nhan_id"` // Chứa mảng JSON hoặc 1 ID
	TieuDe       string        `json:"tieu_de"`
	NoiDung      string        `json:"noi_dung"`
	DinhKem      []FileDinhKem `json:"dinh_kem"`
	ThamChieuID  string        `json:"tham_chieu_id"`
	ReplyChoID   string        `json:"reply_cho_id"`
	NgayTao      string        `json:"ngay_tao"`
	NguoiDoc     []string      `json:"nguoi_doc"`
	TrangThaiXoa []string      `json:"trang_thai_xoa"`

	TenNguoiGui    string `json:"ten_nguoi_gui,omitempty"`
	ChucVuNguoiGui string `json:"chuc_vu_nguoi_gui,omitempty"`
	AvatarNguoiGui string `json:"avatar_nguoi_gui,omitempty"`
	DaDoc          bool   `json:"da_doc"`
}

var (
	CacheTinNhan = make(map[string][]*TinNhan)
	mtxTinNhan   sync.RWMutex
)

func NapTinNhan(shopID string) {
	if shopID == "" { shopID = cau_hinh.BienCauHinh.IdFileSheet }
	raw, err := LoadSheetData(shopID, "TIN_NHAN")
	if err != nil { return }

	list := []*TinNhan{}
	for i, r := range raw {
		if i < DongBatDau_TinNhan-1 { continue }
		id := LayString(r, CotTN_MaTinNhan)
		if id == "" { continue }

		tn := &TinNhan{
			SpreadsheetID:  shopID,
			DongTrongSheet: i + 1,
			MaTinNhan:      id,
			LoaiTinNhan:    LayString(r, CotTN_LoaiTinNhan),
			NguoiGuiID:     LayString(r, CotTN_NguoiGuiID),
			NguoiNhanID:    LayString(r, CotTN_NguoiNhanID),
			TieuDe:         LayString(r, CotTN_TieuDe),
			NoiDung:        LayString(r, CotTN_NoiDung),
			ThamChieuID:    LayString(r, CotTN_ThamChieuID),
			ReplyChoID:     LayString(r, CotTN_ReplyChoID),
			NgayTao:        LayString(r, CotTN_NgayTao),
		}

		_ = json.Unmarshal([]byte(LayString(r, CotTN_DinhKemJson)), &tn.DinhKem)
		if tn.DinhKem == nil { tn.DinhKem = make([]FileDinhKem, 0) }

		_ = json.Unmarshal([]byte(LayString(r, CotTN_NguoiDocJson)), &tn.NguoiDoc)
		if tn.NguoiDoc == nil { tn.NguoiDoc = make([]string, 0) }

		_ = json.Unmarshal([]byte(LayString(r, CotTN_TrangThaiXoa)), &tn.TrangThaiXoa)
		if tn.TrangThaiXoa == nil { tn.TrangThaiXoa = make([]string, 0) }

		list = append(list, tn)
	}

	mtxTinNhan.Lock()
	CacheTinNhan[shopID] = list
	mtxTinNhan.Unlock()
}

func ContainsString(slice []string, val string) bool {
	for _, item := range slice {
		if item == val { return true }
	}
	return false
}

func LayHopThuNguoiDung(shopID string, maKH string, vaiTro string) []*TinNhan {
	mtxTinNhan.RLock()
	defer mtxTinNhan.RUnlock()
	
	allMsgs := CacheTinNhan[shopID]
	var inbox []*TinNhan
	
	for _, m := range allMsgs {
		if ContainsString(m.TrangThaiXoa, maKH) { continue }

		isReceiver := false
		
		// [THUẬT TOÁN MỚI]: Đọc mảng JSON nếu là tin nhắn ALL
		if m.LoaiTinNhan == "ALL" {
			var listNhan []string
			if err := json.Unmarshal([]byte(m.NguoiNhanID), &listNhan); err == nil {
				if ContainsString(listNhan, maKH) {
					isReceiver = true
				}
			}
		} else if m.NguoiNhanID == maKH {
			isReceiver = true
		} else if strings.HasPrefix(m.NguoiNhanID, "ROLE_") {
			roleTarget := strings.TrimPrefix(m.NguoiNhanID, "ROLE_")
			if vaiTro == roleTarget { isReceiver = true }
		}
		
		if isReceiver {
			msgCopy := *m 
			msgCopy.DaDoc = ContainsString(m.NguoiDoc, maKH)
			
			if m.NguoiGuiID == "SYSTEM" || m.NguoiGuiID == "0000000000000000000" {
				// Lấy info thực của Bot nếu có, nếu không thì Hardcode
				bot, ok := LayKhachHang(shopID, "0000000000000000000")
				if ok {
					msgCopy.TenNguoiGui = bot.TenKhachHang
					msgCopy.ChucVuNguoiGui = bot.ChucVu
				} else {
					msgCopy.TenNguoiGui = "Trợ lý ảo 99K"
					msgCopy.ChucVuNguoiGui = "Hệ thống"
				}
				msgCopy.AvatarNguoiGui = "99"
			} else {
				sender, ok := LayKhachHang(shopID, m.NguoiGuiID)
				if ok {
					msgCopy.TenNguoiGui = sender.TenKhachHang
					msgCopy.ChucVuNguoiGui = sender.ChucVu
					if len(sender.TenKhachHang) > 0 {
						msgCopy.AvatarNguoiGui = string([]rune(sender.TenKhachHang)[0])
					}
				} else {
					msgCopy.TenNguoiGui = "Ẩn danh"
				}
			}
			inbox = append(inbox, &msgCopy)
		}
	}
	return inbox
}

func DanhDauDocTinNhan(shopID string, maKH string, msgID string) {
	mtxTinNhan.Lock()
	defer mtxTinNhan.Unlock()
	
	for _, m := range CacheTinNhan[shopID] {
		if m.MaTinNhan == msgID {
			if !ContainsString(m.NguoiDoc, maKH) {
				m.NguoiDoc = append(m.NguoiDoc, maKH)
				ThemVaoHangCho(shopID, "TIN_NHAN", m.DongTrongSheet, CotTN_NguoiDocJson, ToJSON(m.NguoiDoc))
			}
			break
		}
	}
}

func ThemMoiTinNhan(shopID string, msg *TinNhan) {
	mtxTinNhan.Lock()
	list := CacheTinNhan[shopID]
	
	maxRow := DongBatDau_TinNhan - 1
	for _, m := range list {
		if m.DongTrongSheet > maxRow {
			maxRow = m.DongTrongSheet
		}
	}
	msg.DongTrongSheet = maxRow + 1
	
	CacheTinNhan[shopID] = append(list, msg)
	mtxTinNhan.Unlock()
	
	ghi := ThemVaoHangCho
	sh := "TIN_NHAN"
	r := msg.DongTrongSheet
	
	ghi(shopID, sh, r, CotTN_MaTinNhan, msg.MaTinNhan)
	ghi(shopID, sh, r, CotTN_LoaiTinNhan, msg.LoaiTinNhan)
	ghi(shopID, sh, r, CotTN_NguoiGuiID, msg.NguoiGuiID)
	ghi(shopID, sh, r, CotTN_NguoiNhanID, msg.NguoiNhanID)
	ghi(shopID, sh, r, CotTN_TieuDe, msg.TieuDe)
	ghi(shopID, sh, r, CotTN_NoiDung, msg.NoiDung)
	ghi(shopID, sh, r, CotTN_DinhKemJson, ToJSON(msg.DinhKem))
	ghi(shopID, sh, r, CotTN_ThamChieuID, msg.ThamChieuID)
	ghi(shopID, sh, r, CotTN_ReplyChoID, msg.ReplyChoID)
	ghi(shopID, sh, r, CotTN_NgayTao, msg.NgayTao)
	ghi(shopID, sh, r, CotTN_NguoiDocJson, ToJSON(msg.NguoiDoc))
	ghi(shopID, sh, r, CotTN_TrangThaiXoa, ToJSON(msg.TrangThaiXoa))
}
