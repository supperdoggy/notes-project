$.get( "http://localhost:8080/api/ping", function( data ) {
    $("#result").html( data["result"] + $.cookie("t"));
});