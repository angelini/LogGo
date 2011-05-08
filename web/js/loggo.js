var Loggo = {
    search: function(query) {
        $.ajax({
            type: 'POST',
            url: './search',
            data: {
                search: query
            },
            success: function(data) {
                console.log(data)
            } 
        });
    }
}
