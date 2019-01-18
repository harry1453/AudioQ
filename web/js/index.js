let currentCue = 0;
let numberOfCues = 0;

function getProject() {
    fetch("../api/getProject").then(http => {
        return http.json();
    }).then(project => {
        currentCue = project.CurrentCue;
        numberOfCues = project.Cues.length;

        document.getElementById("projectName").innerText = "Project: " + project.Name;

        document.getElementById("projectNameInput").value = project.Name;
        document.getElementById("bufferSize").value = project.Settings.BufferSize;

        document.getElementById("cues").innerHTML =" <tr><th>#</th><th>Sel</th><th>Cue Name</th><th>Jump</th><th>Rename</th><th>Remove</th></tr>";
        for (let i = 0; i < project.Cues.length; i++) {
            let cue = project.Cues[i];
            let sel = project.CurrentCue === i ? "Sel" : "";

            // TODO HTML escaping
            document.getElementById("cues").innerHTML += "<tr><td>"+i+"</td><td>"+sel+"</td><td>"+cue.Name+"</td><td><button onclick='jumpToCue("+i+");'>Jump!</button></td><td><button onclick='renameCue("+i+", \""+cue.Name+"\");'>Rename</button></td><td><button onclick='removeCueFromProject("+i+");' >Remove</button></td></tr>";
        }
    });
}

function playNext() {
    fetch("../api/playNext", {method: "POST"}).then(http => {
        return http.json();
    }).then(result => {
        if (!result.OK) {
            logError("PlayNext()", result.Error);
        }
        getProject();
    });
}

function stopPlaying() {
    fetch("../api/stopPlaying", {method: "POST"}).then(http => {
        return http.json();
    }).then(result => {
        if (!result.OK) {
            logError("StopPlaying()", result.Error);
        }
        getProject();
    });
}

function moveCue() {
    let from = prompt("From?");
    let to = prompt("To?");
    fetch("../api/moveCue/"+from+"/"+to, {method: "POST"}).then(http => {
        return http.json();
    }).then(result => {
        if (!result.OK) {
            logError("JumpToCue("+cueNumber+")", result.Error);
        }
        getProject();
    });
}

function jumpToCue(cueNumber) {
    fetch("../api/jumpTo/"+cueNumber, {method: "POST"}).then(http => {
        return http.json();
    }).then(result => {
        if (!result.OK) {
            logError("JumpToCue("+cueNumber+")", result.Error);
        }
        getProject();
    });
}

function updateProjectName() {
    let projectName = document.getElementById("projectNameInput").value;
    fetch("../api/updateProjectName/"+projectName, {method: "POST"}).then(http => {
        return http.json();
    }).then(result => {
        if (!result.OK) {
            logError("UpdateProjectName("+projectName+")", result.Error);
        }
        getProject();
    });
}

function updateProjectSettings() {
    let formData = new FormData(document.getElementById("settingsForm"));
    fetch("../api/updateProjectSettings", {method: "POST", body: formData}).then(http => {
        return http.json();
    }).then(result => {
        if (!result.OK) {
            logError("UpdateProjectSettings()", result.Error);
        }
        getProject();
    });
}

function addCueToProject() {
    let formData = new FormData(document.getElementById("addCueForm"));
    fetch("../api/addCue", {method: "POST", body: formData}).then(http => {
        return http.json();
    }).then(result => {
        if (!result.OK) {
            logError("AddCueToProject()", result.Error);
        }
        getProject();
    });
}

function removeCueFromProject(cueNumber) {
    if (!confirm("Are you sure you want to remove cue "+cueNumber+"?")) return;
    fetch("../api/removeCue/"+cueNumber, {method: "POST"}).then(http => {
        return http.json();
    }).then(result => {
        if (!result.OK) {
            logError("RemoveCueFromProject()", result.Error);
        }
        getProject();
    });
}

function renameCue(cueNumber, previousName) {
    let cueName = prompt("New cue name?", previousName);
    if (cueName === previousName || cueName === "" || cueName === null) return;
    fetch("../api/renameCue/"+cueNumber+"/"+cueName, {method: "POST"}).then(http => {
        return http.json();
    }).then(result => {
        if (!result.OK) {
            logError("RenameCue("+cueNumber+", "+cueName+")", result.Error);
        }
        getProject();
    });
}

function loadProject() {
    let formData = new FormData(document.getElementById("loadProjectForm"));
    fetch("../api/loadProject", {method: "POST", body: formData}).then(http => {
        return http.json();
    }).then(result => {
        if (!result.OK) {
            logError("LoadProject()", result.Error);
        }
        getProject();
    });
}

function forwardOne() {
    let newCue = currentCue + 1;
    if (newCue === numberOfCues) {
        newCue = 0;
    }
    jumpToCue(newCue);
}

function backOne() {
    let newCue = currentCue - 1;
    if (newCue === -1) {
        newCue = numberOfCues - 1;
    }
    jumpToCue(newCue);
}

function saveProject() {
    window.location.href = "../api/saveProject";
}

function logError(wasDoing, error) {
    console.log("Error: " + error + ", encountered while doing: " + wasDoing)
}

getProject();

document.onkeypress = function (e) {
    e = e || window.event;
    // use e.keyCode
    switch (e.key) {
        case " ":
            playNext();
            break;

        case "s":
            stopPlaying();
            break;

        case "d":
            forwardOne();
            break;

        case "a":
            backOne();
            break;

        case "w":
            getProject();
            break;
    }
};
