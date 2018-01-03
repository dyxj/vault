document.addEventListener('DOMContentLoaded', function () {
    Particles.init({
        selector: '.background',
        maxParticles: 150,
        sizeVariations: 5,
    });

    var eBtn = document.getElementById("enbtn");
    var dBtn = document.getElementById("debtn");
    var ekey1 = document.getElementById("enKey1");
    var ekey2 = document.getElementById("enKey2");
    var dkey1 = document.getElementById("deKey1");
    var eFile = document.getElementById("enFile");
    var dFile = document.getElementById("deFile");

    eBtn.addEventListener('click', function () {
        const form = new FormData();
        form.append("cryptKey1", ekey1.value);
        form.append("cryptKey2", ekey2.value);
        form.append("usrfile", eFile.files[0]);
        postEncrypt(form)
    });

    dBtn.addEventListener('click', function () {
        const form = new FormData();
        form.append("cryptKey1", dkey1.value);
        form.append("usrfile", dFile.files[0]);
        postDecrypt(form)
    });
});

function postEncrypt(data) {
    axios({
        method: 'post',
        url: '/encrypt',
        data: data,
        responseType: 'arraybuffer'
    }).then(function (response) {
        const url = window.URL.createObjectURL(new Blob([response.data]));
        const link = document.createElement('a');
        link.href = url;
        link.setAttribute('download', response.headers["x-file-name"]);
        document.body.appendChild(link);
        link.click();
    }).catch(function (error) {
        var enc = new TextDecoder();        
        var errmsg = enc.decode(error.response.data) ;
        console.log(errmsg);
        callAlert(errmsg);
    });
}

function postDecrypt(data) {
    axios({
        method: 'post',
        url: '/decrypt',
        data: data,
        responseType: 'arraybuffer'
    }).then(function (response) {
        const url = window.URL.createObjectURL(new Blob([response.data]));
        const link = document.createElement('a');
        link.href = url;
        link.setAttribute('download', response.headers["x-file-name"]);
        document.body.appendChild(link);
        link.click();
    }).catch(function (error) {
        var enc = new TextDecoder();
        var errmsg = enc.decode(error.response.data) ;
        console.log(errmsg);
        callAlert(errmsg);
    });
}

function callAlert(text) {
    document.getElementById("error-alert").innerHTML = text;
    showAlert();
    setTimeout(hideAlert, 2000);
}

function showAlert() {
    $("#error-alert")
        .css('opacity', 0)
        .slideDown('slow')
        .animate(
        { opacity: 1 },
        { queue: false, duration: 'slow' }
        );
}

function hideAlert() {
    console.log("hide");
    $("#error-alert")
        .css('opacity', 1)
        .slideUp('slow')
        .animate(
        { opacity: 0 },
        { queue: false, duration: 'slow' }
        );
}

// handle custom file inputs
$('.custom-file-input').on('change', function () {
    let fileName = $(this).val().split('\\').pop();
    $(this).next('.custom-file-label').html(fileName);
});