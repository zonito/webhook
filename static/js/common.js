var wh = {
  util: {
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
  },
  ajax: {
    handler: {
      board: function (data) {
        $('#boards').html('<option value="0">-- Choose board --</option>');
        $.each(data, function (index) {
          var temp = data[index];
          if (index === 0) {
            $('#board_name').val(temp.name);
          }
          $('#boards').append(
            '<option value="' + temp.id + '">' + temp.name + '</option>');
        });
      },
      list: function (data) {
        $('#lists').html('');
        $.each(data, function (index) {
          var temp = data[index];
          if (index === 0) {
            $('#list_name').val(temp.name);
          }
          $('#lists').append(
            '<option value="' + temp.id + '">' + temp.name + '</option>');
        });
      },
      hookResponse: function (data) {
        if (!data.success) {
          $('.error').text(data.reason);
        } else {
          $('#toast').attr('text', 'Webhook Created. webhook.co/' + data.handler);
          $('#toast')[0].show();
          $('#modal')[0].toggle();
          setTimeout(function () {
            $('iron-ajax').attr('url', '/created.json?' + data.handler);
          }, 2000);
        }
      },
      showCreatedList: function (data) {
        console.log("success");
      }
    },
    request: {
      board: function () {
        $.ajax({
          method: "POST",
          url: "/trello/boards/list"
        }).done(wh.ajax.handler.board);
      },
      list: function (value) {
        $.ajax({
          method: 'POST',
          url: '/trello/lists/' + value
        }).done(wh.ajax.handler.list);
      },
      createHook: function (data) {
        $.ajax({
          url: '/save',
          method: 'POST',
          data: data
        }).done(wh.ajax.handler.hookResponse);
      },
      createdList: function () {
        $.ajax({
          url: '/created.json',
          type: 'GET'
        }).done(wh.ajax.handler.showCreatedList);
      }
    }
  },
  event: {
    handler: {
      boards: function () {
        $('#board_name').val($("option:selected", this).text());
        $('.list').show();
        var value = $(this).val();
        if (!value.length || value === "0") {
          $('#lists').html('');
          $('.list').hide();
        }
        wh.ajax.request.list(value);
      },
      lists: function () {
        $('[name="list_name"]').val($("option:selected", this).text());
      },
      openCreateDialog: function () {
        $('#modal')[0].toggle();
      },
      showForm: function () {
        var service = $(this).attr('alt');
        if (service === 'trello' && !$('#boards').text().trim().length) {
          wh.ajax.request.board();
        }
        $('.hide').hide();
        $('img.selected').removeClass('selected');
        $('.' + service).show();
        $('.buttons .hide').show();
        if (!$('.buttons').hasClass('pad')) {
          $('.buttons').addClass('pad');
        }
        $('.error').text('');
        $(this).addClass('selected');
      },
      createHook: function () {
        var service = $('img.selected').attr('alt');
        var data = {
          'service': service
        };
        if (service === 'trello') {
          data.board_name = $('#board_name').val();
          data.board_id = $('#boards').val();
          data.list_name = $('#list_name').val();
          data.list_id = $('#lists').val();
          if (!data.board_name.length || !data.board_id.length || !data.list_id.length || !data.list_name.length) {
            $('.error').text('Provide all information.');
            return;
          }
        } else if (service === 'telegram') {
          data.tele_code = $('#tele_code').val();
          if (data.tele_code.length !== 6) {
            $('.error').text('Invalid code.');
            return;
          }
        }
        wh.ajax.request.createHook(data);
      }
    },
    add: function () {
      var handler = wh.event.handler;
      $('#lists').on('change', handler.lists);
      $('#boards').on('change', handler.boards);
      $('#addHook').on('click', handler.openCreateDialog);
      $('.app img').on('click', handler.showForm);
      $('#create').on('click', handler.createHook);
    }
  }
};

$(document).ready(function () {
  wh.event.add();
});
