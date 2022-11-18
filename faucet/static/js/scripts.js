
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
        loader();

        $.ajax({
            type: "POST",
            url: "http://localhost:8000/fund",
            data: data,
            success: function(resp) {
                successAlert();
            },
            error: function(xhr) {
                resp = $.parseJSON(xhr.responseText);
                errorAlert(resp.errors[0]);
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

function loader(){
    $('#result-msg').html(`
<div class="spinner-grow text-light" role="status">
  <span class="sr-only"></span>
</div>`);
}
