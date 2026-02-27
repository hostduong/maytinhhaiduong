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
	
	// [ĐÃ FIX]: Thêm biến ảo để Front-end biết tin này đã đọc hay chưa
	DaDoc          bool          `json:"da_doc"` 
}

var (
	CacheTinNhan = make(map[string][]*TinNhan)
	mtxTinNhan   sync.RWMutex
)

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

func LayHopThuNguoiDung(shopID string, maKH string, vaiTro string) []*TinNhan {
	mtxTinNhan.RLock()
	defer mtxTinNhan.RUnlock()

	var inbox []*TinNhan
	for _, m := range CacheTinNhan[shopID] {
		if ContainsString(m.TrangThaiXoa, maKH) { continue }

		isReceiver := false
		
		if m.LoaiTinNhan == "CHAT" || m.LoaiTinNhan == "AUTO" {
			if m.NguoiNhanID == maKH || m.NguoiGuiID == maKH {
				isReceiver = true
			}
		} else if m.LoaiTinNhan == "SYSTEM" || m.LoaiTinNhan == "ALL" {
			if strings.Contains(m.NguoiNhanID, "["+maKH+"]") || strings.Contains(m.NguoiNhanID, "\""+maKH+"\"") {
				isReceiver = true
			} else if m.NguoiNhanID == "ALL" || strings.Contains(m.NguoiNhanID, vaiTro) {
				isReceiver = true
			} else if m.NguoiGuiID == maKH {
				isReceiver = true 
			}
		}

		// [ĐÃ FIX]: Tạo bản Copy để nhét cờ DaDoc riêng cho từng User, không làm hỏng Cache tổng
		if isReceiver { 
			mCopy := *m
			mCopy.DaDoc = ContainsString(m.NguoiDoc, maKH)
			inbox = append(inbox, &mCopy) 
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
		if m.DongTrongSheet > maxRow { maxRow = m.DongTrongSheet }
	}
	msg.DongTrongSheet = maxRow + 1
	
	CacheTinNhan[shopID] = append(list, msg)
	mtxTinNhan.Unlock()
	
	if msg.DinhKem == nil { msg.DinhKem = make([]FileDinhKem, 0) }
	if msg.NguoiDoc == nil { msg.NguoiDoc = make([]string, 0) }
	if msg.TrangThaiXoa == nil { msg.TrangThaiXoa = make([]string, 0) }

	dongMoi := []interface{}{
		msg.MaTinNhan,        
		msg.LoaiTinNhan,      
		msg.NguoiGuiID,       
		msg.NguoiNhanID,      
		msg.TieuDe,           
		msg.NoiDung,          
		ToJSON(msg.DinhKem),  
		msg.ThamChieuID,      
		msg.ReplyChoID,       
		msg.NgayTao,          
		ToJSON(msg.NguoiDoc), 
		ToJSON(msg.TrangThaiXoa), 
	}

	ThemDongVaoHangCho(shopID, "TIN_NHAN", dongMoi)
}
