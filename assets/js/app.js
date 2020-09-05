function onResponse(response) {
    var divNotification = document.querySelector("#notice");
    var divNotificationMessage = document.querySelector("#notice-message");
    var className = (response.status !== 400) ? "alert-success" : "alert-danger";
    divNotificationMessage.innerHTML = response.data;

    divNotification.classList.add(className);

    document.getElementById("notice").style.display = "block";
    setTimeout(function () {
        document.getElementById("notice").style.display = "none";
        divNotification.classList.remove(className);
    },3000);
}

function senddata() {
    var file = document.getElementById("inputFile").files[0];

    if (!file) {
        return;
    }
    upload(file);

    function upload(file) {
        var formData = new FormData();
        formData.append('file', file);

        post('/apps/uploadReport', formData)
            .then(onResponse)
            .catch(onResponse)
    }
}(document, axios);

function addapp() {
    var name = document.getElementById("name").value;
    var url = document.getElementById("url").value;
    var language = document.getElementById("language").value
    var framework = document.getElementById("framework").value

    uploadAppData(name, url, language, framework);

    name.reset();
    url.reset();
    language.reset();
    framework.reset();

    function uploadAppData(name,url,language,framework) {
        var formData = new FormData();
        formData.append('name', name);
        formData.append('url', url);
        formData.append('language', language);
        formData.append('framework', framework);

        post('/apps/addApp', formData)
            .then(onResponse)
            .catch(onResponse)


    }
}(document, axios);

