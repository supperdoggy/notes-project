function deleteUser(id){
    $.ajax("http://localhost:2283/api/deleteUser",{
        type:"POST",
        data:{"id":id,},
    })
    window.alert("deleted!")
    window.location.href = "http://localhost:8080/admin"
}