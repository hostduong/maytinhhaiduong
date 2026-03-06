// ===================================================================
// XỬ LÝ NGHIỆP VỤ MUA GÓI TẠI TRANG CHỦ ADMIN
// ===================================================================

// 1. Khai báo các hàm Alert dùng chung (kế thừa từ phong cách Master)
function showSuccessAlert(msg) { 
    return Swal.fire({ 
        title: '<div class="flex flex-col items-center gap-2"><div class="w-16 h-16 bg-purple-50 text-purple-600 rounded-full flex items-center justify-center mb-2"><svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7"></path></svg></div><span class="text-3xl font-black title-gradient-purple">Thành Công</span></div>', 
        html: `<div class="text-[15px] font-bold text-slate-600 mt-2">${msg}</div>`, 
        customClass: { popup: 'swal-master-success' }, 
        showConfirmButton: false, timer: 2000 
    }); 
}

function showErrorAlert(msg) { 
    return Swal.fire({ 
        title: '<div class="flex flex-col items-center gap-2"><div class="w-16 h-16 bg-orange-50 text-orange-600 rounded-full flex items-center justify-center mb-2"><svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path></svg></div><span class="text-3xl font-black text-orange-600">Thất Bại</span></div>', 
        html: `<div class="text-[15px] font-bold text-orange-900 mt-2">${msg}</div>`, 
        customClass: { popup: 'swal-master-error' }, 
        confirmButtonText: 'Tôi hiểu rồi', buttonsStyling: false, confirmButtonClass: 'bg-gradient-to-r from-orange-500 to-red-500 text-white px-8 py-3 rounded-xl font-bold uppercase tracking-widest text-xs shadow-md' 
    }); 
}

// 2. Gọi API để Check giá Real-time khi nhập mã Code (Dùng chung API với bang-gia)
async function checkGiaTuServer(maGoi, maCode) {
    const fd = new FormData();
    fd.append('ma_goi', maGoi);
    fd.append('ma_code', maCode);
    const res = await fetch('/bang-gia/api/check-price', { method: 'POST', body: fd });
    return await res.json();
}

// 3. Hàm bật Modal và xử lý luồng mua gói
async function moModalMua(maGoi, tenGoi, giaGoc) {
    const fPrice = new Intl.NumberFormat('vi-VN').format(giaGoc);
    
    Swal.fire({
        title: `<span class="text-2xl font-black title-gradient-purple-swal">Gói ${tenGoi}</span>`,
        html: `
            <p class="text-sm font-bold text-slate-600 mb-2">Giá: ${fPrice}₫</p>
            <div class="relative mt-4">
                <input type="text" id="inp_code" class="input-premium font-mono !text-xs text-center border-slate-200" placeholder="NHẬP MÃ GIẢM GIÁ (NẾU CÓ)">
            </div>
            <div id="price_box" class="text-center p-3 mt-4 bg-purple-50 rounded-xl border-2 border-transparent transition-colors">
                <p class="text-[10px] font-bold text-purple-600 uppercase tracking-widest">Giá Cuối Cùng</p>
                <p class="text-3xl font-black text-gradient-animated" id="final_txt">${fPrice}₫</p>
            </div>
        `,
        showCancelButton: true, confirmButtonText: 'Xác nhận mua', cancelButtonText: 'Hủy',
        buttonsStyling: false,
        customClass: { 
            popup: 'swal-animated-purple', 
            confirmButton: 'btn-premium px-8 py-3 rounded-xl mx-2 shadow-lg shadow-purple-200', 
            cancelButton: 'bg-slate-100 text-slate-500 px-8 py-3 rounded-xl mx-2 font-bold text-xs uppercase tracking-widest hover:bg-slate-200 transition' 
        },
        didOpen: () => {
            const inp = document.getElementById('inp_code');
            const txt = document.getElementById('final_txt');
            const box = document.getElementById('price_box');

            inp.addEventListener('input', async (e) => {
                const code = e.target.value.trim();
                if(code.length < 2) {
                    txt.innerText = fPrice + '₫';
                    box.classList.remove('border-emerald-500');
                    return;
                }
                const data = await checkGiaTuServer(maGoi, code);
                if(data.status === 'ok') {
                    txt.innerText = new Intl.NumberFormat('vi-VN').format(data.final_price) + '₫';
                    if(data.is_valid) { box.classList.add('border-emerald-500'); } else { box.classList.remove('border-emerald-500'); }
                }
            });
        }
    }).then(async (result) => {
        if (result.isConfirmed) {
            const code = document.getElementById('inp_code').value.trim();
            Swal.fire({ title: 'Đang xử lý giao dịch...', allowOutsideClick: false, didOpen: () => Swal.showLoading() });

            const fd = new FormData();
            fd.append('ma_goi', maGoi);
            fd.append('ma_code', code);

            try {
                // Tận dụng luôn API Mua Gói đang có sẵn của hệ thống
                const res = await fetch('/bang-gia/api/mua-goi', { method: 'POST', body: fd });
                const data = await res.json();

                if(data.status === 'ok') { 
                    showSuccessAlert("Đăng ký thành công! Hệ thống đang khởi tạo...").then(() => {
                        window.location.href = data.redirect_url; 
                    });
                } else { 
                    showErrorAlert(data.msg); 
                }
            } catch(e) {
                showErrorAlert("Lỗi đường truyền! Vui lòng thử lại.");
            }
        }
    });
}
