function getProject() {
    fetch("../api/getProject").then(http => {
        return http.json();
    }).then(project => {
        document.getElementById("projectTitle").innerText = "Project: " + project.Name;

        document.getElementById("bufferSize").innerText = "Buffer Size: " + project.Settings.BufferSize;

        document.getElementById("cues").innerHTML =" <tr><th>#</th><th>Sel</th><th>Cue Name</th><th>Jump</th></tr>";
        for (let i = 0; i < project.Cues.length; i++) {
            let cue = project.Cues[i];
            let sel = project.CurrentCue === i ? "Sel" : "";
            document.getElementById("cues").innerHTML += "<tr><td>"+i+"</td><td>"+sel+"</td><td>"+cue.Name+"</td><td><button onclick='jumpToCue("+i+")'>Jump!</button></td></tr>";
        }
    });
}

function playNext() {
    fetch("../api/playNext", {method: "POST"}).then(http => {
        return http.json();
    }).then(result => {
        if (!result.OK) {
            logError("PlayNext()");
        }
        getProject();
    })
}

function jumpToCue(cueNumber) {
    fetch("../api/jumpTo/"+cueNumber, {method: "POST"}).then(http => {
        return http.json();
    }).then(result => {
        if (!result.OK) {
            logError("JumpToCue("+cueNumber+")");
        }
        getProject();
    })
}

function logError(wasDoing, error) {
    console.log("Error: " + error + ", encountered while doing: " + wasDoing)
}

getProject();
