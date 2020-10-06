"use strict"; // ES6
window.onload = () => {

    let http = {
        post: (path, data) => {
            return new Promise((resolve, reject) => {
                let xhr = new XMLHttpRequest();
                xhr.open("POST", path, true);
                if (JSON.stringify(data).indexOf("json-content-type") > -1) {
                    xhr.setRequestHeader("Content-Type", "application/json")
                }
                xhr.onreadystatechange = () => {
                    if (xhr.readyState === XMLHttpRequest.DONE) return resolve(xhr);
                };
                xhr.send(data);
            });
        }
    };

    let ui = {
        output: document.getElementById("output"),
        image: document.querySelector("img#img"),
        btnFile: document.getElementById("by-file"),
        btnBase64: document.getElementById("by-base64"),
        cancel: document.getElementById("cancel-input"),
        file: document.getElementById("file"),
        langs: document.querySelector("input[name=langs]"),
        whitelist: document.querySelector("input[name=whitelist]"),
        hocr: document.querySelector("input[name=hocr]"),
        submit: document.getElementById("submit"),
        loading: document.querySelector("button#submit>span:first-child"),
        standby: document.querySelector("button#submit>span:last-child"),
        show: uri => ui.image.setAttribute("src", uri),
        clear: () => {
            ui.image.setAttribute("src", ""), ui.file.value = '';
        },
        start: () => {
            ui.loading.style.display = "block";
            ui.standby.style.display = "none";
            ui.submit.setAttribute("disabled", true);
            ui.output.innerText = "{}";
        },
        finish: () => {
            ui.loading.style.display = "none";
            ui.standby.style.display = "block";
            ui.submit.removeAttribute("disabled");
        },
    };

    ui.file.addEventListener("change", ev => {
        if (!ev.target.files || !ev.target.files.length) return null;
        const r = new FileReader();
        r.onload = e => ui.show(e.target.result);
        r.readAsDataURL(ev.target.files[0]);
    });
    ui.btnFile.addEventListener("click", () => ui.file.click());
    ui.btnBase64.addEventListener("click", () => {
        const uri = window.prompt("Please paste your base64 image URI");
        if (uri) {
            ui.clear();
            ui.show(uri);
        }
    });
    ui.cancel.addEventListener("click", () => {
        ui.clear()
        alert("eewewe")
    });
    ui.submit.addEventListener("click", () => {
        ui.start();
        const req = generateRequest();
        if (!req) return ui.finish();
        http.post(req.path, req.data).then(xhr => {
            ui.output.innerText = `${xhr.status} ${xhr.statusText}\n-----\n${xhr.response}`;
            ui.finish();
        }).catch(() => ui.finish());
    })

    let generateRequest = () => {
        let req = {path: "", data: null};
        if (ui.file.files && ui.file.files.length != 0) {
            req.path = "/api/ocr/file";
            req.data = new FormData();
            if (ui.langs.value) req.data.append("languages", ui.langs.value);
            if (ui.whitelist.value) req.data.append("whitelist", ui.whitelist.value);
            if (ui.hocr.checked) req.data.append("hocrMode", true);
            req.data.append("file", ui.file.files[0]);
        } else if (/^data:.+/.test(ui.image.src)) {
            req.path = "/api/ocr/base64";
            let data = {base64: ui.image.src, "json-content-type": true};
            if (ui.langs.value) data["languages"] = ui.langs.value;
            if (ui.whitelist.value) data["whitelist"] = ui.whitelist.value;
            if (ui.hocr.checked) data["hocrMode"] = true;
            req.data = JSON.stringify(data);
            console.log('dataaa', req.data)
        } else {
            return window.alert("no image input set");
        }
        return req;
    };
};

let addPixel = () => {
    let lastDiv = getLastDivObj()

    let nextDivId = getNextDivId(1)
    let newDiv = document.createElement(nextDivId)
}

let getNextDivId = (index) => {
    const prefix = "pixel-group"
    let div = document.getElementById(prefix + index)

    if (!div)
        return prefix + index

    getNextDivId(index + 1)
}

let getLastDivId = (index) => {
    const prefix = "pixel-group"
    let div = document.getElementById("" + prefix + index)
    if (!div) {
        return "" + prefix + (index - 1)
    }

    return getNextDivId(index + 1)
}

let getLastDivObj = () => {
    let lastDivId = getLastDivId(1)
    console.log('结果 lastDivId', lastDivId)

    return document.getElementById(lastDivId)
}