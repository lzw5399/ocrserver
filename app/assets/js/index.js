"use strict"; // ES6
const containerId = "pixel-container"
const pixelPrefix = "pixel-group"
let existItems = []

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
    ui.cancel.addEventListener("click", () => ui.clear());
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
        removeErrorIfExist()
        let req = {path: "", data: null};
        if (ui.file.files && ui.file.files.length !== 0 && isEmpty()) {
            req.path = "/api/ocr/file";
            req.data = new FormData();
            if (ui.langs.value) req.data.append("languages", ui.langs.value);
            if (ui.whitelist.value) req.data.append("whitelist", ui.whitelist.value);
            if (ui.hocr.checked) req.data.append("hocrMode", true);
            req.data.append("file", ui.file.files[0]);
        } else if (ui.file.files && ui.file.files.length !== 0 && !isEmpty()) {
            req.path = "/api/ocr/scan-crop-file";
            req.data = new FormData();
            if (isFinalValid()) {
                req.data.append("matrixPixels", genMatrixPixels());
            } else {
                return showPixelErrorMsg()
            }
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
        } else {
            return window.alert("no image input set");
        }
        return req;
    };

    addPixel(true)
    setInterval(() => {
        let container = document.getElementById(containerId)
        let tempArr = []
        container.childNodes.forEach((v, i, p) => {
            if (v && v.id && v.id.indexOf(pixelPrefix) > -1) {
                tempArr.push(getIndexFromDivId(v.id))
            }
        })
        existItems = tempArr
    }, 500)
    setInterval(() => {
        removeErrorIfExist()
    }, 5000)
};

let addPixel = (initial) => {
    // 创建新的div
    let nextDivId = getNextDivId()
    let nextDiv = document.createElement("div")
    nextDiv.id = nextDivId

    // 创建新div下的元素
    let nextIndex = getNextIndex()
    let pointA = createStrong("点A &nbsp;")
    let pointB = createStrong("点B &nbsp;")
    let ax = createNumberInput("ax" + nextIndex, "x")
    let ay = createNumberInput("ay" + nextIndex, "y")
    let bx = createNumberInput("bx" + nextIndex, "x")
    let by = createNumberInput("by" + nextIndex, "y")
    let iconAdd = createIconElem("add" + nextIndex, nextIndex, true)
    let iconMinus = createIconElem("minus" + nextIndex, nextIndex, false)
    let spaceArray = createSpaceArray(5)

    nextDiv.appendChild(pointA)
    nextDiv.appendChild(ax)
    nextDiv.appendChild(spaceArray[0])
    nextDiv.appendChild(ay)
    nextDiv.appendChild(spaceArray[1])
    nextDiv.appendChild(pointB)
    nextDiv.appendChild(bx)
    nextDiv.appendChild(spaceArray[2])
    nextDiv.appendChild(by)
    nextDiv.appendChild(spaceArray[3])
    nextDiv.appendChild(iconAdd)
    nextDiv.appendChild(spaceArray[4])
    if (!initial) {
        nextDiv.appendChild(iconMinus)
    }

    // 最终将新div添加到旧div下
    let container = document.getElementById(containerId)
    container.appendChild(nextDiv)

    existItems.push(nextIndex)
    refreshIconShown()
}

let removePixel = (index) => {
    let container = document.getElementById(containerId)
    let currentDiv = document.getElementById(pixelPrefix + index)
    container.removeChild(currentDiv)

    for (let i = 0; i < existItems.length; i++) {
        if (existItems[i] === index) {
            existItems = existItems.splice(i, 1)
        }
    }

    refreshIconShown()
}

let createNumberInput = (id, placeholder) => {
    let elem = document.createElement("input")
    elem.setAttribute("class", "pixel input")
    elem.setAttribute("type", "number")
    elem.setAttribute("id", id)
    elem.setAttribute("placeholder", placeholder)
    elem.setAttribute("style", "margin:1px;")

    return elem
}

let createStrong = (text) => {
    let elem = document.createElement("strong")
    elem.innerHTML = text

    return elem
}

let createIconElem = (id, index, isAdd) => {
    let elem = document.createElement("i")
    elem.setAttribute("id", id)
    elem.setAttribute("style", "cursor:pointer")

    if (isAdd) {
        elem.setAttribute("class", "fas fa-plus")
        elem.setAttribute("onclick", "addPixel()")
    } else {
        elem.setAttribute("class", "fas fa-minus")
        elem.setAttribute("onclick", `removePixel(${index})`)
    }

    return elem
}

let createSpaceArray = (num) => {
    let arr = []
    for (let i = 1; i <= num; i++) {
        let elem = document.createElement("span")
        elem.innerHTML = "&nbsp;"
        arr.push(elem)
    }

    return arr
}

// 获取next不代表是最新的
let getNextDivId = () => {
    let nextIndex = getNextIndex()

    return pixelPrefix + nextIndex
}

// 获取next不代表是最新的
let getNextIndex = () => {
    let lastDivId = getLastDivId(1)
    let lastIndex = lastDivId.replace(pixelPrefix, "")

    return (Number(lastIndex) + 1)
}

// 这个也不是存在中最后面的
let getLastDivId = (index) => {
    let div = document.getElementById("" + pixelPrefix + index)
    if (!div) {
        return "" + pixelPrefix + (index - 1)
    }

    return getLastDivId(index + 1)
}

let genMatrixPixels = () => {
    let finalArr = []

    for (let i = 0; i < existItems.length; i++) {
        finalArr.push(getPointsByIndex(existItems[i]))
    }

    return JSON.stringify(finalArr)
}

let isEmpty = () => {
    if (existItems.length !== 1)
        return false

    let index = existItems[0]
    let ax = document.getElementById("ax" + index)
    let ay = document.getElementById("ay" + index)
    let bx = document.getElementById("bx" + index)
    let by = document.getElementById("by" + index)

    return !ax.value && !ay.value && !bx.value && !by.value
}

let isFinalValid = () => {
    let valid = true
    for (let i = 0; i < existItems.length; i++) {
        let groupValid = isValidGroup(existItems[i])
        if (!groupValid) {
            valid = false
            break
        }
    }

    return valid
}

let isValidGroup = (index) => {
    let ax = document.getElementById("ax" + index)
    let ay = document.getElementById("ay" + index)
    let bx = document.getElementById("bx" + index)
    let by = document.getElementById("by" + index)

    return ax.value && ay.value && bx.value && by.value
}

let getPointsByIndex = (index) => {
    let ax = document.getElementById("ax" + index)
    let ay = document.getElementById("ay" + index)
    let bx = document.getElementById("bx" + index)
    let by = document.getElementById("by" + index)

    return {
        pointA: {
            x: Number(ax.value),
            y: Number(ay.value),
        },
        pointB: {
            x: Number(bx.value),
            y: Number(by.value)
        }
    }
}

let getIndexFromDivId = (divId) => {
    return Number(divId.replace(pixelPrefix, ""))
}

let showPixelErrorMsg = () => {
    removeErrorIfExist()

    let container = document.getElementById(containerId)
    let errMsg = document.createElement("p")
    errMsg.setAttribute("id", "err-msg")
    errMsg.setAttribute("class", "help")
    errMsg.style.color = "red"
    errMsg.innerHTML = "矩阵像素点存在空值，请确认后重试！"
    container.appendChild(errMsg)
}

let removeErrorIfExist = () => {
    let errMsg = document.getElementById("err-msg")
    if (!errMsg) return

    let container = document.getElementById(containerId)
    container.removeChild(errMsg)
}

let refreshIconShown = () => {
    if (existItems.length === 1) {
        removeIconByIndexIfExist(existItems[0], false)
        createIconByIndexIfNotExist(existItems[0], true)
        return
    }

    // 最后一个之前都是减号
    for (let i = 0; i < existItems.length - 1; i++) {
        removeIconByIndexIfExist(existItems[i], true)
        createIconByIndexIfNotExist(existItems[i], false)
    }
    removeIconByIndexIfExist(existItems[existItems.length - 1], false)
    createIconByIndexIfNotExist(existItems[existItems.length - 1], true)
}

let createIconByIndexIfNotExist = (index, isAdd) => {
    let iconId = isAdd ? "add" + index : "minus" + index
    let iconElem = document.getElementById(iconId)

    if (iconElem)
        return

    let group = document.getElementById(pixelPrefix + index)
    let i = createIconElem(iconId, index, isAdd)
    group.appendChild(i)
}

let removeIconByIndexIfExist = (index, isAdd) => {
    let iconId = isAdd ? "add" + index : "minus" + index
    let iconElem = document.getElementById(iconId)

    if (!iconElem)
        return

    let group = document.getElementById(pixelPrefix + index)
    group.removeChild(iconElem)
}
