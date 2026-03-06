const isMeRoot = (window.MasterConfig.currentUserID === "0000000000000000001");
const myLevel = window.MasterConfig.myLevel;
let itiPhone, itiZalo;

function checkInputState(input) {
    if(input.value.trim() !== "") input.classList.add('has-data');
    else input.classList.remove('has-data');
}

function formatSocialLink(val, type) {
    val = val.trim().replace(/\s/g, ''); 
    if (val === "") return "";
    let username = val;
    if (type === 'facebook') { username = username.replace(/^(https?:\/\/)?(www\.)?facebook\.com\//i, '').replace(/^(https?:\/\/)?fb\.com\//i, ''); } 
    else if (type === 'tiktok') { username = username.replace(/^(https?:\/\/)?(www\.)?tiktok\.com\//i, '').replace(/^(https?:\/\/)?vt\.tiktok\.com\//i, ''); }
    username = username.split('?')[0].replace(/\/+$/, '');
    if (username.startsWith('@')) username = username.substring(1);
    if (username === "") return val;
    if (type === 'facebook') return `https://www.facebook.com/${username}/`;
    if (type === 'tiktok') return `https://www.tiktok.com/@${username}`; 
    return val;
}

document.addEventListener("DOMContentLoaded", () => {
    document.querySelectorAll('.input-premium').forEach(input => {
        checkInputState(input); 
        input.addEventListener('input', function() {
            if (this.type !== 'email' && this.id !== 'f_fb' && this.id !== 'f_tiktok' && this.id !== 'f_avatar' && this.type !== 'password') {
                this.value = this.value.replace(/^[=+\-@]/, '').replace(/[<>;'`{}\[\]\\]/g, '');
            }
            checkInputState(this); 
        });
    });

    if (document.getElementById('f_fb')) document.getElementById('f_fb').addEventListener('blur', function() { this.value = formatSocialLink(this.value, 'facebook'); checkInputState(this); });
    if (document.getElementById('f_tiktok')) document.getElementById('f_tiktok').addEventListener('blur', function() { this.value = formatSocialLink(this.value, 'tiktok'); checkInputState(this); });
    if (document.getElementById('f_avatar')) {
        document.getElementById('f_avatar').addEventListener('blur', function() { if (this.value.trim() !== "" && !/^https?:\/\//i.test(this.value.trim())) this.value = 'https://' + this.value.trim(); checkInputState(this); });
        document.getElementById('f_avatar').addEventListener('input', function() { if(/\s/.test(this.value)) this.value = this.value.replace(/\s/g, ''); });
    }
});

function togglePass(id) {
    const input = document.getElementById(id);
    const eyeSlash = document.getElementById('eye_slash_' + id); const eyeOpen = document.getElementById('eye_open_' + id);
    if (input.type === "password") { input.type = "text"; eyeSlash.classList.add('hidden'); eyeOpen.classList.remove('hidden'); } 
    else { input.type = "password"; eyeOpen.classList.add('hidden'); eyeSlash.classList.remove('hidden'); }
}

window.toggleSwalPin = function() {
    const input = document.getElementById('swal-pin-input');
    const eyeSlash = document.getElementById('eye_slash_swal'); const eyeOpen = document.getElementById('eye_open_swal');
    if (input.type === "password") { input.type = "text"; eyeSlash.classList.add('hidden'); eyeOpen.classList.remove('hidden'); } 
    else { input.type = "password"; eyeOpen.classList.add('hidden'); eyeSlash.classList.remove('hidden'); }
};

const mapData = {};
document.querySelectorAll('.m-data').forEach(el => {
    let id = el.getAttribute('data-id');
    mapData[id] = {
        id: id, ten: el.getAttribute('data-name'), user: el.getAttribute('data-user'), email: el.getAttribute('data-email'),
        role: el.getAttribute('data-role'), title: el.getAttribute('data-title'), sdt: el.getAttribute('data-phone'), dob: el.getAttribute('data-dob'), gender: el.getAttribute('data-gender'),
        address: el.getAttribute('data-address'), tax: el.getAttribute('data-tax'), status: el.getAttribute('data-status'),
        zalo: el.getAttribute('data-zalo'), fb: el.getAttribute('data-fb'), tiktok: el.getAttribute('data-tiktok'), note: el.getAttribute('data-note'),
        avatar: el.getAttribute('data-avatar'), source: el.getAttribute('data-source'), created: el.getAttribute('data-created'), updater: el.getAttribute('data-updater'), updated: el.getAttribute('data-updated')
    };
});

function filterTable(inputId, tbodyId) { 
    const queryClean = (str => str ? str.normalize('NFD').replace(/[\u0300-\u036f]/g, '') : '')(document.getElementById(inputId).value.toLowerCase().trim());
    document.querySelectorAll(`#${tbodyId} tr.hover-soft-row`).forEach(row => { 
        const cb = row.querySelector('.chk-member'); if (!cb) return; const data = mapData[cb.value]; if (!data) return;
        let statusText = data.status == 1 ? "hoat dong" : (data.status == -1 ? "doi xoa" : "khoa tam khoa");
        const combined = `${data.id} ${data.ten} ${data.user} ${data.email} ${data.role} ${data.title} ${data.sdt} ${data.dob} ${data.address} ${data.tax} ${data.zalo} ${data.fb} ${data.tiktok} ${data.source} ${data.note} ${statusText}`.toLowerCase();
        if ((str => str ? str.normalize('NFD').replace(/[\u0300-\u036f]/g, '') : '')(combined).includes(queryClean)) { row.style.display = ''; } else { row.style.display = 'none'; }
    }); 
}

function closeModal(id) { document.getElementById(id).classList.add('hidden'); }

function editMember(ma) {
    let data = mapData[ma]; if(!data) return;
    document.getElementById('formEdit').reset();
    document.getElementById('f_pass').setAttribute('readonly', 'true'); document.getElementById('f_pin').setAttribute('readonly', 'true');
    
    if (!itiPhone) itiPhone = window.intlTelInput(document.querySelector("#f_sdt_fake"), { initialCountry: "vn", separateDialCode: true, utilsScript: "https://cdn.jsdelivr.net/npm/intl-tel-input@18.2.1/build/js/utils.js" });
    if (!itiZalo) itiZalo = window.intlTelInput(document.querySelector("#f_zalo_fake"), { initialCountry: "vn", separateDialCode: true, utilsScript: "https://cdn.jsdelivr.net/npm/intl-tel-input@18.2.1/build/js/utils.js" });
    
    document.getElementById('f_ma').value = ma; document.getElementById('f_ten').value = data.ten;
    itiPhone.setNumber(data.sdt || ""); itiZalo.setNumber(data.zalo || ""); 
    document.getElementById('f_dob').value = data.dob; document.getElementById('f_gender').value = data.gender;
    document.getElementById('f_tax').value = data.tax; document.getElementById('f_address').value = data.address;
    document.getElementById('f_title').value = data.title; document.getElementById('f_avatar').value = data.avatar; 
    document.getElementById('f_source').value = data.source; document.getElementById('f_fb').value = data.fb; 
    document.getElementById('f_tiktok').value = data.tiktok; document.getElementById('f_note').value = data.note; 
    document.getElementById('f_log_created').innerText = data.created || "--"; document.getElementById('f_log_updated').innerText = data.updated || "--"; document.getElementById('f_log_updater').innerText = data.updater || "Hệ Thống";

    const isSystemBot = (ma === "0000000000000000000");
    if (isSystemBot) { document.getElementById('box_cap_lai_bao_mat').classList.add('hidden'); } else { document.getElementById('box_cap_lai_bao_mat').classList.remove('hidden'); }

    let roleSelect = document.getElementById('f_role'); let statusSelect = document.getElementById('f_status');
    Array.from(roleSelect.options).forEach(opt => { opt.disabled = (opt.value === "quan_tri_he_thong" && ma !== "0000000000000000001"); });

    roleSelect.value = data.role; statusSelect.value = data.status; roleSelect.disabled = false; statusSelect.disabled = false;
    if (ma === "0000000000000000001" && window.MasterConfig.currentUserID !== "0000000000000000001") { roleSelect.disabled = true; statusSelect.disabled = true; } 
    if (ma === window.MasterConfig.currentUserID) { statusSelect.disabled = true; }

    document.querySelectorAll('.input-premium').forEach(input => checkInputState(input));
    document.getElementById('modalEdit').classList.remove('hidden');
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
    if (getSelectedIDs().length === 0) { 
        Swal.fire({
            title: '<div class="flex flex-col items-center gap-3"><div class="w-16 h-16 bg-gradient-to-br from-red-100 to-orange-100 text-orange-600 rounded-full flex items-center justify-center shadow-inner border border-orange-200"><svg class="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" /></svg></div><span class="text-3xl font-black title-gradient-red">Cảnh Báo!</span></div>',
            html: `<div class="text-[15px] text-orange-900 font-medium px-4">Vui lòng tích chọn ít nhất 1 người nhận.</div>`,
            customClass: { popup: 'swal-master-error', confirmButton: 'bg-gradient-to-r from-orange-500 to-red-500 text-white font-bold rounded-xl px-8 py-3 shadow-lg' }, buttonsStyling: false
        }); return; 
    } 
    document.getElementById('msgTitle').value = ""; document.getElementById('msgContent').value = ""; updateCount(); document.getElementById('modalMsg').classList.remove('hidden'); 
}

async function saveMember() {
    const { value: pinXacNhan } = await Swal.fire({
        title: '<div class="flex flex-col items-center gap-3"><div class="w-14 h-14 bg-gradient-to-br from-purple-100 to-indigo-100 text-purple-600 rounded-full flex items-center justify-center shadow-inner border border-purple-200"><svg class="w-7 h-7" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"></path></svg></div><span class="text-3xl font-black title-gradient-purple">Xác thực Mã PIN</span></div>',
        customClass: { popup: 'swal-master-success', confirmButton: 'bg-gradient-to-r from-purple-600 to-indigo-600 text-white font-bold rounded-xl px-8 py-3 shadow-lg', cancelButton: 'bg-slate-100 text-slate-600 font-bold rounded-xl px-8 py-3 border border-slate-200' },
        buttonsStyling: false, 
        html: `
            <div class="relative w-full max-w-[280px] mx-auto group mt-4">
                <input type="password" id="swal-pin-input" class="relative w-full bg-white border-2 border-slate-200 rounded-2xl py-4 px-4 text-center tracking-[0.25em] font-black text-2xl outline-none focus:border-purple-500" maxlength="8" autocomplete="new-password" readonly onfocus="this.removeAttribute('readonly');" placeholder="••••••••" oninput="this.value = this.value.replace(/[^0-9]/g, ''); this.classList.remove('pin-error'); document.getElementById('swal-custom-err').classList.add('hidden'); checkPinLength(this);">
            </div>
            <div id="swal-custom-err" class="hidden text-red-500 mt-4 text-[13px] font-bold">Vui lòng nhập đủ 8 số!</div>
        `,
        showCancelButton: true, confirmButtonText: 'Xác nhận Lưu', cancelButtonText: 'Hủy bỏ', reverseButtons: true, 
        didOpen: () => { setTimeout(() => { const input = document.getElementById('swal-pin-input'); input.removeAttribute('readonly'); input.focus(); }, 100); },
        preConfirm: () => { 
            const input = document.getElementById('swal-pin-input'); const errBox = document.getElementById('swal-custom-err'); 
            if (input.value.length < 8) { input.classList.add('pin-error'); errBox.classList.remove('hidden'); return false; } 
            return input.value; 
        }
    });

    if (!pinXacNhan) return;
    const btn = document.getElementById('btnSaveMem'); const old = btn.innerHTML;
    btn.disabled = true; btn.innerHTML = 'Đang xử lý...';
    
    const fd = new FormData(document.getElementById('formEdit'));
    fd.append('pin_xac_nhan', pinXacNhan);
    if (itiPhone) fd.set('dien_thoai', itiPhone.getNumber());
    if (itiZalo) fd.set('zalo', itiZalo.getNumber());
    if(document.getElementById('f_role').disabled) { fd.append('vai_tro', mapData[document.getElementById('f_ma').value].role); }
    if(document.getElementById('f_status').disabled) { fd.append('trang_thai', mapData[document.getElementById('f_ma').value].status); }

    try {
        const res = await fetch('/master/api/thanh-vien/save', { method: 'POST', body: fd });
        const data = await res.json();
        if(data.status === 'ok') {
            Swal.fire({ icon: 'success', title: 'Thành Công', text: data.msg, showConfirmButton: false, timer: 1500 }).then(() => location.reload());
        } else {
            Swal.fire({ icon: 'error', title: 'Từ Chối', text: data.msg }); 
        }
    } catch(e) { Swal.fire({ icon: 'error', title: 'Lỗi', text: 'Mất kết nối máy chủ!' }); } 
    finally { btn.disabled = false; btn.innerHTML = old; }
}

async function sendMessages() {
    let title = document.getElementById('msgTitle').value.trim(); let content = document.getElementById('msgContent').value.trim();
    if (!title || !content) { Swal.fire({ icon: 'error', text: 'Tiêu đề và nội dung không được bỏ trống!' }); return; }
    
    const btn = document.getElementById('btnSendMsg'); const old = btn.innerHTML; 
    btn.disabled = true; btn.innerHTML = 'Đang phát sóng...';
    
    const fd = new FormData(); fd.append('tieu_de', title); fd.append('noi_dung', content); fd.append('danh_sach_id', JSON.stringify(getSelectedIDs())); 
    let chkBot = document.getElementById('chkSendAsBot'); if (chkBot && chkBot.checked) fd.append('send_as_bot', "1");

    try {
        const res = await fetch('/master/api/thanh-vien/send-msg', { method: 'POST', body: fd });
        const data = await res.json();
        if(data.status === 'ok') { Swal.fire({ icon: 'success', text: data.msg, timer: 2000 }).then(() => location.reload()); } 
        else { Swal.fire({ icon: 'error', text: data.msg }); }
    } catch(e) { Swal.fire({ icon: 'error', text: 'Mất kết nối máy chủ!' }); } 
    finally { btn.disabled = false; btn.innerHTML = old; }
}
