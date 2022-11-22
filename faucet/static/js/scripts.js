const FAUCET_BACKEND="{{.}}";
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
            url: FAUCET_BACKEND,
            crossDomain: true, // set as a cross domain request
            data: data,
            timeout: 60_000,
            success: function(data, status, xhr) {
                successAlert();
            },
            error: function(jqXhr, textStatus, errorThrown) {
                console.log("ajax error: ", errorThrown)
                if (jqXhr != null && jqXhr.responseText != null ) {
                    resp = $.parseJSON(jqXhr.responseText);
                    errorAlert(resp.errors[0]);
                } else {
                    errorAlert(errorThrown);
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

function loader(){
    $('#result-msg').html(`
<div class="spinner-grow text-light" role="status">
  <span class="sr-only"></span>
</div>`);
}
