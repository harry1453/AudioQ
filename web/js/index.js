function getProject() {
    fetch("../api/getProject").then(http => {
        return http.json();
    }).then(project => {
        console.log(project);
    });
}

getProject();
