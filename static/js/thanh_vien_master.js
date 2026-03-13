const currentUserID = window.MasterConfig.currentUserID; 
const isMeRoot = (currentUserID === "0000000000000000001");
const myLevel = window.MasterConfig.myLevel;
let itiPhone, itiZalo;

// ==========================================
// CƠ CHẾ KÉO THẢ (DRAG TO RESIZE SPLIT PANE)
// ==========================================
const listPane = document.getElementById('listPaneUI');
const detailPane = document.getElementById('modalEdit'); // Panel bên phải
const dragBar = document.getElementById('dragResizerUI');
const splitContainer = document.getElementById('splitUIContainer');
let isResizing = false;

if (dragBar) {
    dragBar.addEventListener('mousedown', function(e) {
        isResizing = true;
        document.body.classList.add('cursor-col-resize', 'select-none');
        listPane.style.transition = 'none'; 
    });

    document.addEventListener('mousemove', function(e) {
        if (!isResizing) return;
        e.preventDefault(); 
        let offsetLeft = e.clientX - splitContainer.getBoundingClientRect().left;
        let percentage = (offsetLeft / splitContainer.getBoundingClientRect().width) * 100;
        if (percentage < 30) percentage = 30; 
        if (percentage > 70) percentage = 70; 
        listPane.style.width = percentage + '%';
        listPane.style.flex = 'none'; 
        detailPane.style.width = (100 - percentage) + '%';
    });

    document.addEventListener('mouseup', function(e) {
        if (isResizing) {
            isResizing = false;
            document.body.classList.remove('cursor-col-resize', 'select-none');
            listPane.style.transition = ''; 
        }
    });
}

function syncTopAlignment() {
    if(window.innerWidth >= 1280) { 
        const leftBox = document.querySelector('.animated-border-wrapper');
        const spacer = document.getElementById('formSpacer');
        if(leftBox && splitContainer && spacer) {
            const topOffset = leftBox.getBoundingClientRect().top - splitContainer.getBoundingClientRect().top;
            spacer.style.height = topOffset + 'px';
        }
    } else {
        const spacer = document.getElementById('formSpacer');
        if(spacer) spacer.style.height = '0px';
    }
}

// ==========================================
// NGHIỆP VỤ XỬ LÝ DỮ LIỆU
// ==========================================
function checkInputState(input) {
    if(input.value.trim() !== "") input.classList.add('has-data');
    else input.classList.remove('has-data');
}

function formatSocialLink(val, type) {
    val = val.trim().replace(/\s/g, ''); 
    if (val === "") return "";
    let username = val;
    if (type === 'facebook') {
        username = username.replace(/^(https?:\/\/)?(www\.)?facebook\.com\//i, '');
        username = username.replace(/^(https?:\/\/)?fb\.com\//i, '');
    } else if (type === 'tiktok') {
        username = username.replace(/^(https?:\/\/)?(www\.)?tiktok\.com\//i, '');
        username = username.replace(/^(https?:\/\/)?vt\.tiktok\.com\//i, '');
    }
    username = username.split('?')[0];
    username = username.replace(/\/+$/, '');
    if (username.startsWith('@')) username = username.substring(1);
    if (username === "") return val;
    if (type === 'facebook') { return `https://www.facebook.com/${username}/`; } 
    else if (type === 'tiktok') { return `https://www.tiktok.com/@${username}`; }
    return val;
}

document.addEventListener("DOMContentLoaded", () => {
    document.querySelectorAll('.input-premium').forEach(input => {
        checkInputState(input); 
        input.addEventListener('input', function() {
            if (this.type !== 'email' && this.id !== 'f_fb' && this.id !== 'f_tiktok' && this.id !== 'f_avatar' && this.type !== 'password') {
                this.value = this.value.replace(/^[=+\-@]/, '');
                this.value = this.value.replace(/[<>;'`{}\[\]\\]/g, '');
            }
            checkInputState(this); 
        });
    });

    const fbInput = document.getElementById('f_fb');
    if (fbInput) { fbInput.addEventListener('blur', function() { this.value = formatSocialLink(this.value, 'facebook'); checkInputState(this); }); }
    
    const tiktokInput = document.getElementById('f_tiktok');
    if (tiktokInput) { tiktokInput.addEventListener('blur', function() { this.value = formatSocialLink(this.value, 'tiktok'); checkInputState(this); }); }
    
    const avatarInput = document.getElementById('f_avatar');
    if (avatarInput) {
        avatarInput.addEventListener('blur', function() {
            let v = this.value.trim();
            if (v !== "" && !/^https?:\/\//i.test(v)) this.value = 'https://' + v;
            checkInputState(this);
        });
        avatarInput.addEventListener('input', function() { if(/\s/.test(this.value)) this.value = this.value.replace(/\s/g, ''); });
    }
});

function togglePass(id) {
    const input = document.getElementById(id);
    if(!input) return;
    if (input.type === "password") { input.type = "text"; } else { input.type = "password"; }
}

window.toggleSwalPin = function() {
    const input = document.getElementById('swal-pin-input');
    const eyeSlash = document.getElementById('eye_slash_swal');
    const eyeOpen = document.getElementById('eye_open_swal');
    if (input.type === "password") { input.type = "text"; eyeSlash.classList.add('hidden'); eyeOpen.classList.remove('hidden'); } 
    else { input.type = "password"; eyeOpen.classList.add('hidden'); eyeSlash.classList.remove('hidden'); }
};

const mapData = {};
document.querySelectorAll('.m-data').forEach(el => {
    let id = el.getAttribute('data-id');
    mapData[id] = {
        id: id, ten: el.getAttribute('data-name'), user: el.getAttribute('data-user'), email: el.getAttribute('data-email'),
        role: el.getAttribute('data-role'), title: el.getAttribute('data-title'), 
        sdt: el.getAttribute('data-phone'), dob: el.getAttribute('data-dob'), gender: el.getAttribute('data-gender'), 
        address: el.getAttribute('data-address'), tax: el.getAttribute('data-tax'), status: el.getAttribute('data-status'),
        zalo: el.getAttribute('data-zalo'), fb: el.getAttribute('data-fb'), tiktok: el.getAttribute('data-tiktok'), 
        note: el.getAttribute('data-note'), avatar: el.getAttribute('data-avatar'), source: el.getAttribute('data-source'),
        created: el.getAttribute('data-created'), updater: el.getAttribute('data-updater'), updated: el.getAttribute('data-updated')
    };
});

function filterTable(inputId, tbodyId) { 
    const query = document.getElementById(inputId).value.toLowerCase().trim(); 
    const removeAccents = (str) => str ? str.normalize('NFD').replace(/[\u0300-\u036f]/g, '') : '';
    const queryClean = removeAccents(query);

    document.querySelectorAll(`#${tbodyId} tr.sp-row`).forEach(row => { 
        const cb = row.querySelector('.chk-member'); if (!cb) return;
        const data = mapData[cb.value]; if (!data) return;

        let statusText = "khoa tam khoa";
        if(data.status == 1) statusText = "hoat dong"; else if (data.status == -1) statusText = "doi xoa";

        const combinedString = `${data.id} ${data.ten} ${data.user} ${data.email} ${data.role} ${data.title} ${data.sdt} ${data.dob} ${data.address} ${data.tax} ${data.zalo} ${data.fb} ${data.tiktok} ${data.source} ${data.note} ${statusText}`.toLowerCase();
        if (removeAccents(combinedString).includes(queryClean)) { row.style.display = ''; } else { row.style.display = 'none'; }
    }); 
}

function closeModal(id) { 
    let modal = document.getElementById(id);
    if (modal) {
        modal.classList.add('hidden');
        modal.classList.remove('is-open');
    }
    if (id === 'modalEdit') {
        let db = document.getElementById('dragResizerUI');
        if(db) db.classList.remove('is-open');
        
        let lp = document.getElementById('listPaneUI');
        if(lp) lp.style.width = '100%';
        
        let dp = document.getElementById('modalEdit');
        if(dp) dp.style.width = '';
        
        document.querySelectorAll('.sp-row').forEach(r => r.classList.remove('row-active-purple'));
    }
}

function editMember(ma) {
    let data = mapData[ma]; if(!data) return;
    
    document.getElementById('formEdit').reset();
    document.getElementById('f_pass').setAttribute('readonly', 'true'); document.getElementById('f_pin').setAttribute('readonly', 'true');
    
    const inputPhone = document.querySelector("#f_sdt_fake"); const inputZalo = document.querySelector("#f_zalo_fake"); 
    if (!itiPhone) itiPhone = window.intlTelInput(inputPhone, { initialCountry: "vn", separateDialCode: true, utilsScript: "https://cdn.jsdelivr.net/npm/intl-tel-input@18.2.1/build/js/utils.js" });
    if (!itiZalo) itiZalo = window.intlTelInput(inputZalo, { initialCountry: "vn", separateDialCode: true, utilsScript: "https://cdn.jsdelivr.net/npm/intl-tel-input@18.2.1/build/js/utils.js" });
    
    // Đổ động Tên thành viên lên thay cho chữ "Hồ Sơ & Phân Quyền"
    let titleEl = document.getElementById('modalTitle');
    if (titleEl) { titleEl.innerHTML = `Sửa: <span class="text-purple-600 ml-1 font-mono tracking-wider">${data.ten}</span>`; }

    document.getElementById('f_ma').value = ma; 
    if (document.getElementById('f_ma_hien_thi')) document.getElementById('f_ma_hien_thi').value = ma;
    if (document.getElementById('f_user')) document.getElementById('f_user').value = data.user || "";

    document.getElementById('f_ten').value = data.ten;
    itiPhone.setNumber(data.sdt || ""); itiZalo.setNumber(data.zalo || ""); 
    document.getElementById('f_dob').value = data.dob; document.getElementById('f_gender').value = data.gender;
    document.getElementById('f_tax').value = data.tax; document.getElementById('f_address').value = data.address;
    document.getElementById('f_title').value = data.title; document.getElementById('f_avatar').value = data.avatar; 
    document.getElementById('f_source').value = data.source; document.getElementById('f_fb').value = data.fb; 
    document.getElementById('f_tiktok').value = data.tiktok; document.getElementById('f_note').value = data.note; 
    
    const formatDate = (ts) => {
        if (!ts || ts === "0") return "--";
        const date = new Date(parseInt(ts) * 1000);
        return date.toLocaleString('vi-VN', {day:'2-digit', month:'2-digit', year:'numeric', hour:'2-digit', minute:'2-digit'});
    };
    document.getElementById('f_log_created').innerText = formatDate(data.created); 
    document.getElementById('f_log_updated').innerText = formatDate(data.updated); 

    const isSystemBot = (ma === "0000000000000000000");
    const boxBaoMat = document.getElementById('box_cap_lai_bao_mat');
    if (isSystemBot) { boxBaoMat.classList.add('hidden'); } else { boxBaoMat.classList.remove('hidden'); }

    let roleSelect = document.getElementById('f_role');
    let statusSelect = document.getElementById('f_status');

    Array.from(roleSelect.options).forEach(opt => { opt.disabled = false; });
    roleSelect.value = data.role; statusSelect.value = data.status;
    roleSelect.disabled = false; statusSelect.disabled = false;

    if (ma === "0000000000000000001") { roleSelect.disabled = true; statusSelect.disabled = true; } 
    if (ma === currentUserID) { statusSelect.disabled = true; }

    document.querySelectorAll('.input-premium').forEach(input => checkInputState(input));
    
    // Đổi màu dòng đang chọn
    document.querySelectorAll('.sp-row').forEach(r => r.classList.remove('row-active-purple'));
    const activeRow = document.querySelector(`.sp-row[data-id="${ma}"]`);
    if(activeRow) activeRow.classList.add('row-active-purple');

    // Mở Pane bẻ đôi màn hình
    let modal = document.getElementById('modalEdit');
    if (modal) {
        modal.classList.remove('hidden');
        modal.classList.add('is-open');
    }
    
    let db = document.getElementById('dragResizerUI');
    if(db) db.classList.add('is-open');

    syncTopAlignment();
}

window.toggleAllCheckboxes = function() { 
    let isChecked = document.getElementById('checkAll').checked; 
    document.querySelectorAll('.chk-member').forEach(cb => cb.checked = isChecked); 
    updateCount(); 
};

window.checkPinLength = function(input) {
    if (input.value.length === 8) {
        input.classList.remove('border-purple-500', 'focus:border-purple-500', 'focus:shadow-[0_0_0_4px_rgba(147,51,234,0.1)]');
        input.classList.add('border-emerald-500', 'focus:border-emerald-500', 'text-emerald-600', 'focus:shadow-[0_0_0_4px_rgba(16,185,129,0.1)]');
        input.style.transform = 'scale(1.03)'; setTimeout(() => input.style.transform = 'scale(1)', 150);
    } else {
        input.classList.remove('border-emerald-500', 'focus:border-emerald-500', 'text-emerald-600', 'focus:shadow-[0_0_0_4px_rgba(16,185,129,0.1)]');
        input.classList.add('border-purple-500', 'focus:border-purple-500', 'text-slate-800', 'focus:shadow-[0_0_0_4px_rgba(147,51,234,0.1)]');
    }
}

function getSelectedIDs() { let ids = []; document.querySelectorAll('.chk-member:checked').forEach(cb => ids.push(cb.value)); return ids; }
function updateCount() { let count = document.querySelectorAll('.chk-member:checked').length; document.getElementById('countSelected').innerText = count; document.getElementById('lblMsgCount').innerText = count; }

function openMessageModal() { 
    let count = getSelectedIDs().length; 
    if (count === 0) { 
        Swal.fire({
            title: '<div class="flex flex-col items-center gap-3"><div class="w-16 h-16 bg-gradient-to-br from-red-100 to-orange-100 text-orange-600 rounded-full flex items-center justify-center shadow-inner border border-orange-200"><svg class="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" /></svg></div><span class="text-3xl font-black title-gradient-red">Cảnh Báo!</span></div>',
            html: `<div class="text-[15px] text-orange-900 font-medium px-4">Vui lòng tích chọn ít nhất 1 người nhận ở bảng danh sách.</div>`,
            customClass: { popup: 'swal-master-error', confirmButton: 'bg-gradient-to-r from-orange-500 to-red-500 hover:from-orange-600 hover:to-red-600 text-white font-bold rounded-xl px-8 py-3 shadow-lg shadow-orange-200 transition transform hover:-translate-y-0.5 active:scale-95' },
            buttonsStyling: false, confirmButtonText: 'Đã hiểu'
        });
        return; 
    } 
    document.getElementById('msgTitle').value = ""; document.getElementById('msgContent').value = ""; updateCount(); document.getElementById('modalMsg').classList.remove('hidden'); 
}

async function saveMember() {
    const form = document.getElementById('formEdit');
    if (!form.checkValidity()) { form.reportValidity(); return; }

    const { value: pinXacNhan } = await Swal.fire({
        title: '<div class="flex flex-col items-center gap-3"><div class="w-14 h-14 bg-gradient-to-br from-purple-100 to-indigo-100 text-purple-600 rounded-full flex items-center justify-center shadow-inner border border-purple-200"><svg class="w-7 h-7" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"></path></svg></div><span class="text-3xl font-black title-gradient-purple">Bảo mật 2 lớp</span></div>',
        customClass: { popup: 'swal-master-success', confirmButton: 'bg-gradient-to-r from-purple-600 to-indigo-600 hover:from-purple-700 hover:to-indigo-700 text-white font-bold rounded-xl px-8 py-3 shadow-lg shadow-purple-200 transition transform hover:-translate-y-0.5 active:scale-95', cancelButton: 'bg-slate-100 hover:bg-slate-200 text-slate-600 font-bold rounded-xl px-8 py-3 transition active:scale-95 border border-slate-200' },
        buttonsStyling: false, 
        html: `
            <div class="text-[14px] text-slate-600 mb-8 font-medium leading-relaxed">Hệ thống yêu cầu xác minh danh tính.<br>Vui lòng nhập <b class="text-purple-700">Mã PIN</b> của bạn.</div>
            <div class="relative w-full max-w-[280px] mx-auto group">
                <div class="absolute inset-0 bg-purple-500/20 blur-xl rounded-full opacity-0 group-focus-within:opacity-100 transition duration-500 pointer-events-none"></div>
                <input type="password" id="swal-pin-input" class="relative w-full box-border bg-white border-2 border-slate-200 rounded-2xl py-4 px-4 text-center tracking-[0.25em] font-black text-2xl text-slate-800 outline-none transition-all duration-300 focus:border-purple-500 focus:shadow-[0_0_0_4px_rgba(147,51,234,0.1)] placeholder-slate-300" maxlength="8" autocomplete="new-password" readonly onfocus="this.removeAttribute('readonly');" placeholder="••••••••" oninput="this.value = this.value.replace(/[^0-9]/g, ''); this.classList.remove('pin-error'); document.getElementById('swal-custom-err').classList.add('hidden'); checkPinLength(this);">
                <button type="button" onclick="toggleSwalPin()" class="absolute inset-y-0 right-2 w-10 flex items-center justify-center text-slate-400 hover:text-purple-600 transition z-10 bg-white my-1 rounded-xl">
                    <svg id="eye_slash_swal" class="w-6 h-6 block" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M3.98 8.223A10.477 10.477 0 001.934 12C3.226 16.338 7.244 19.5 12 19.5c.993 0 1.953-.138 2.863-.395M6.228 6.228A10.45 10.45 0 0112 4.5c4.756 0 8.773 3.162 10.065 7.498a10.523 10.523 0 01-4.293 5.774M6.228 6.228L3 3m3.228 3.228l3.65 3.65m7.894 7.894L21 21m-3.228-3.228l-3.65-3.65m0 0a3 3 0 10-4.243-4.243m4.242 4.242L9.88 9.88" /></svg>
                    <svg id="eye_open_swal" class="w-6 h-6 hidden" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M2.036 12.322a1.012 1.012 0 010-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178z" /><path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" /></svg>
                </button>
            </div>
            <div id="swal-custom-err" class="hidden text-red-500 bg-red-50 border border-red-100 rounded-xl px-4 py-2 text-[13px] font-bold mt-4 flex items-center justify-center gap-1.5 w-full max-w-[280px] mx-auto shadow-sm">
                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" /></svg> Vui lòng nhập đủ 8 số!
            </div>
        `,
        showCancelButton: true, confirmButtonText: 'Xác nhận Lưu', cancelButtonText: 'Hủy bỏ', reverseButtons: true, 
        didOpen: () => { setTimeout(() => { const input = document.getElementById('swal-pin-input'); input.removeAttribute('readonly'); input.focus(); }, 100); },
        preConfirm: () => { 
            const input = document.getElementById('swal-pin-input'); const errBox = document.getElementById('swal-custom-err'); const val = input.value; 
            if (val.length < 8) { input.classList.add('pin-error'); errBox.classList.remove('hidden'); errBox.classList.add('animate-bounce'); setTimeout(() => errBox.classList.remove('animate-bounce'), 1000); return false; } 
            return val; 
        }
    });

    if (!pinXacNhan) { return; }

    const btn = document.getElementById('btnSaveMem'); const old = btn.innerHTML;
    btn.disabled = true; btn.classList.add('opacity-80', 'cursor-not-allowed');
    btn.innerHTML = '<div class="flex items-center justify-center gap-2.5"><div class="spinner-led-btn"></div> Đang xử lý...</div>';
    
    const fd = new FormData(form);
    fd.append('pin_xac_nhan', pinXacNhan);
    if (itiPhone) fd.set('dien_thoai', itiPhone.getNumber());
    if (itiZalo) fd.set('zalo', itiZalo.getNumber());
    
    // [LUẬT BẢO VỆ] NẾU Ô BỊ KHÓA, TRUYỀN LẠI DỮ LIỆU CŨ LÊN SERVER ĐỂ KHÔNG BỊ TRỐNG
    if(document.getElementById('f_role').disabled) { fd.append('vai_tro', mapData[document.getElementById('f_ma').value].role); }
    if(document.getElementById('f_status').disabled) { fd.append('trang_thai', mapData[document.getElementById('f_ma').value].status); }

    try {
        const res = await fetch('/master/api/thanh-vien/save', { method: 'POST', body: fd });
        const data = await res.json();
        if(data.status === 'ok') {
            Swal.fire({
                title: '<div class="flex flex-col items-center gap-3"><div class="w-16 h-16 bg-gradient-to-br from-purple-100 to-indigo-100 text-purple-600 rounded-full flex items-center justify-center shadow-inner border border-purple-200"><svg class="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" /></svg></div><span class="text-3xl font-black title-gradient-purple">Thành Công</span></div>',
                html: `<div class="text-[15px] text-slate-600 font-medium px-4">${data.msg}</div>`,
                customClass: { popup: 'swal-master-success' },
                showConfirmButton: false, timer: 1500
            }).then(() => location.reload());
        } else {
            Swal.fire({
                title: '<div class="flex flex-col items-center gap-3"><div class="w-16 h-16 bg-gradient-to-br from-red-100 to-orange-100 text-orange-600 rounded-full flex items-center justify-center shadow-inner border border-orange-200"><svg class="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" /></svg></div><span class="text-3xl font-black title-gradient-red">Từ Chối!</span></div>',
                html: `<div class="text-[15px] text-orange-900 font-medium px-4">${data.msg}</div>`,
                customClass: { popup: 'swal-master-error', confirmButton: 'bg-gradient-to-r from-orange-500 to-red-500 hover:from-orange-600 hover:to-red-600 text-white font-bold rounded-xl px-8 py-3 shadow-lg shadow-orange-200 transition transform hover:-translate-y-0.5 active:scale-95' },
                buttonsStyling: false, confirmButtonText: 'Đã hiểu'
            }); 
        }
    } catch(e) { 
        Swal.fire({
            title: '<div class="flex flex-col items-center gap-3"><div class="w-16 h-16 bg-gradient-to-br from-red-100 to-orange-100 text-orange-600 rounded-full flex items-center justify-center shadow-inner border border-orange-200"><svg class="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" /></svg></div><span class="text-3xl font-black title-gradient-red">Lỗi Kết Nối</span></div>',
            html: `<div class="text-[15px] text-orange-900 font-medium px-4">Đường truyền đến máy chủ bị gián đoạn!</div>`,
            customClass: { popup: 'swal-master-error', confirmButton: 'bg-gradient-to-r from-orange-500 to-red-500 hover:from-orange-600 hover:to-red-600 text-white font-bold rounded-xl px-8 py-3 shadow-lg shadow-orange-200 transition transform hover:-translate-y-0.5 active:scale-95' },
            buttonsStyling: false, confirmButtonText: 'Đóng lại'
        }); 
    } 
    finally { btn.disabled = false; btn.innerHTML = old; btn.classList.remove('opacity-80', 'cursor-not-allowed'); }
}

async function sendMessages() {
    let title = document.getElementById('msgTitle').value.trim(); let content = document.getElementById('msgContent').value.trim(); let ids = getSelectedIDs();
    let chkBot = document.getElementById('chkSendAsBot'); let sendAsBot = (chkBot && chkBot.checked) ? "1" : "0";
    
    if (!title || !content) { 
        Swal.fire({
            title: '<div class="flex flex-col items-center gap-3"><div class="w-16 h-16 bg-gradient-to-br from-red-100 to-orange-100 text-orange-600 rounded-full flex items-center justify-center shadow-inner border border-orange-200"><svg class="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" /></svg></div><span class="text-3xl font-black title-gradient-red">Thiếu Thông Tin</span></div>',
            html: `<div class="text-[15px] text-orange-900 font-medium px-4">Tiêu đề và nội dung không được bỏ trống!</div>`,
            customClass: { popup: 'swal-master-error', confirmButton: 'bg-gradient-to-r from-orange-500 to-red-500 hover:from-orange-600 hover:to-red-600 text-white font-bold rounded-xl px-8 py-3 shadow-lg shadow-orange-200 transition transform hover:-translate-y-0.5 active:scale-95' },
            buttonsStyling: false, confirmButtonText: 'Nhập lại'
        });
        return; 
    }
    
    const btn = document.getElementById('btnSendMsg'); const old = btn.innerHTML; 
    btn.disabled = true; btn.classList.add('opacity-80', 'cursor-not-allowed');
    btn.innerHTML = '<div class="flex items-center justify-center gap-2.5"><div class="spinner-led-btn"></div> Đang phát sóng...</div>';
    
    const fd = new FormData(); fd.append('tieu_de', title); fd.append('noi_dung', content); fd.append('danh_sach_id', JSON.stringify(ids)); fd.append('send_as_bot', sendAsBot);

    try {
        const res = await fetch('/master/api/thanh-vien/send-msg', { method: 'POST', body: fd });
        const data = await res.json();
        if(data.status === 'ok') { 
            Swal.fire({
                title: '<div class="flex flex-col items-center gap-3"><div class="w-16 h-16 bg-gradient-to-br from-purple-100 to-indigo-100 text-purple-600 rounded-full flex items-center justify-center shadow-inner border border-purple-200"><svg class="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" /></svg></div><span class="text-3xl font-black title-gradient-purple">Đã Phát Sóng</span></div>',
                html: `<div class="text-[15px] text-slate-600 font-medium px-4">${data.msg}</div>`,
                customClass: { popup: 'swal-master-success' },
                showConfirmButton: false, timer: 2000
            }).then(() => location.reload()); 
        } 
        else {
            Swal.fire({
                title: '<div class="flex flex-col items-center gap-3"><div class="w-16 h-16 bg-gradient-to-br from-red-100 to-orange-100 text-orange-600 rounded-full flex items-center justify-center shadow-inner border border-orange-200"><svg class="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" /></svg></div><span class="text-3xl font-black title-gradient-red">Lỗi Hệ Thống</span></div>',
                html: `<div class="text-[15px] text-orange-900 font-medium px-4">${data.msg}</div>`,
                customClass: { popup: 'swal-master-error', confirmButton: 'bg-gradient-to-r from-orange-500 to-red-500 hover:from-orange-600 hover:to-red-600 text-white font-bold rounded-xl px-8 py-3 shadow-lg shadow-orange-200 transition transform hover:-translate-y-0.5 active:scale-95' },
                buttonsStyling: false, confirmButtonText: 'Đóng lại'
            });
        }
    } catch(e) { 
        Swal.fire({
            title: '<div class="flex flex-col items-center gap-3"><div class="w-16 h-16 bg-gradient-to-br from-red-100 to-orange-100 text-orange-600 rounded-full flex items-center justify-center shadow-inner border border-orange-200"><svg class="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" /></svg></div><span class="text-3xl font-black title-gradient-red">Lỗi Kết Nối</span></div>',
            html: `<div class="text-[15px] text-orange-900 font-medium px-4">Đường truyền đến máy chủ bị gián đoạn!</div>`,
            customClass: { popup: 'swal-master-error', confirmButton: 'bg-gradient-to-r from-orange-500 to-red-500 hover:from-orange-600 hover:to-red-600 text-white font-bold rounded-xl px-8 py-3 shadow-lg shadow-orange-200 transition transform hover:-translate-y-0.5 active:scale-95' },
            buttonsStyling: false, confirmButtonText: 'Đóng lại'
        }); 
    } 
    finally { btn.disabled = false; btn.innerHTML = old; btn.classList.remove('opacity-80', 'cursor-not-allowed'); }
}
