
// When DOM is loaded this 
// function will get executed
$(() => {
    // function will get executed 
    // on click of submit button
    $('#faucet').on('submit', function(e){
        e.preventDefault();
        ser = $(this).serialize();
        data = JSON.stringify( { address: ser.split("=")[1]} );
        console.log('data', data);

        $.ajax({
            type: "POST",
            url: "http://localhost:8000/fund",
            data: data,
            success: function(resp) {
                alert('success');
                console.log("SUCCESS");
                // TODO: Handle the error and then print the alert.
                console.log(resp);
            },
            error: function() {
                console.log("ERROR");
            }
        });
    });});

