function del(id) {
    $.ajax("http://localhost:2020/api/delete",{
        type:"POST",
        data:{"id":id},
    })
    location.reload();
}

function newNote(username){
    $.ajax("http://localhost:2020/api/newNote", {
        type:"POST",
        data:{"Username":username, "Title":"Example"}
    })
    window.location.href = "http://localhost:8080/"
}