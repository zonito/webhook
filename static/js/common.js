var wh = {
  clip: function (text) {
    var copyElement = document.createElement('input');
    copyElement.setAttribute('type', 'text');
    copyElement.setAttribute('value', text);
    copyElement = document.body.appendChild(copyElement);
    copyElement.select();
    try {
      document.execCommand('copy');
    } finally {
      copyElement.remove();
    }
  }
};

$(document).ready(function () {
  $.ajax({
    method: "POST",
    url: "/trello/boards/list"
  }).done(function (data) {
    $('#boards').html('');
    $.each(data, function (index) {
      var temp = data[index];
      if (index === 0) {
        $('[name="board_name"]').val(temp.name);
      }
      $('#boards').append(
        '<option value="' + temp.id + '">' + temp.name + '</option>');
    });
  });
  $('#lists').on('change', function () {
    $('[name="list_name"]').val($("option:selected", this).text());
  });

  // on change of board
  $('#boards').on('change', function () {
    $('[name="board_name"]').val($("option:selected", this).text());
    var value = $(this).val();
    $.ajax({
      method: 'POST',
      url: '/trello/lists/' + value
    }).done(function (data) {
      $('#lists').html('');
      $.each(data, function (index) {
        var temp = data[index];
        if (index === 0) {
          $('[name="list_name"]').val(temp.name);
        }
        $('#lists').append(
          '<option value="' + temp.id + '">' + temp.name + '</option>');
      });
    });
  });
});
