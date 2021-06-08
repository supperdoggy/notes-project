
function update(id){
    $.ajax("http://localhost:2020/api/updateNote",{
        type:"POST",
        data:{"id":id, "Text":$("#text").val(), "Title":$("#title").val(),
        },
    })
    window.alert("Saved!")
}

function share(id, owner, username) {
    $.ajax("http://localhost:2020/api/share",{
        type:"POST",
        data:{"Id":id, "Username":username, "Owner":owner},
    })
    window.alert("Shared!")
}