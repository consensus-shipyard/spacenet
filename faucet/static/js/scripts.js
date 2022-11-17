
// When DOM is loaded this 
// function will get executed
$(() => {
    // function will get executed 
    // on click of submit button
    $('#faucet').on('submit', function(e){
        e.preventDefault();
        ser = $(this).serialize();
        data = JSON.stringify( { address: ser.split("=")[1]} );
        console.log('request sent:', data);

        $.ajax({
            type: "POST",
            url: "http://localhost:8000/fund",
            data: data,
            success: function(resp) {
                resp = $.parseJSON(resp);
                console.log("response:", resp);
                console.log("error:", resp.Error);
                if (resp.Error == ""){
                    successAlert();
                } else {
                    errorAlert(resp.Error);
                }
            },
            error: function(xhr) {
                if (!xhr.responseText) {
                    errorAlert("unknown error");
                } else {
                    errorAlert(xhr.responseText);
                }
            }
        });
    });});

function successAlert() {
    $('#result-msg').html(`<div class="alert alert-success" role="alert">
  Congratulations! Your Spacenet funds are on their way! 👾
  </div>`);
}

function errorAlert(err) {
    $('#result-msg').html(`<div class="alert alert-danger" role="alert">
  Error requesting token funds: ${err} 🫠
  </div>`);
}
