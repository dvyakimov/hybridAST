function senddata() {
    var file = document.getElementById("inputFile").files[0];

    if (!file) {
        return;
    }
    upload(file);

    function upload(file) {
        var formData = new FormData();
        formData.append('file', file);

        post('/upload', formData)
            .then(onResponse)
            .catch(onResponse)
    }

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
}(document, axios);
