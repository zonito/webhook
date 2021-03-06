/*jslint browser: true */
/*global $,document,window*/

(function () {
  "use strict";

  var wh = {
    util: {
      clip: function (text) {
        var copyElement = document.createElement("input");
        copyElement.setAttribute("type", "text");
        copyElement.setAttribute("value", text);
        copyElement = document.body.appendChild(copyElement);
        copyElement.select();
        try {
          document.execCommand("copy");
        } finally {
          copyElement.remove();
        }
      }
    },
    ajax: {
      handler: {
        board: function (data) {
          $("#boards").html("<option value='0'>-- Choose board --</option>");
          $.each(data, function (index) {
            var temp = data[index];
            if (index === 0) {
              $("#boardName").val(temp.name);
            }
            $("#boards").append(
              "<option value='" + temp.id + "''>" + temp.name + "</option>");
          });
        },
        list: function (data) {
          $("#lists").html("");
          $.each(data, function (index) {
            var temp = data[index];
            if (index === 0) {
              $("#listName").val(temp.name);
            }
            $("#lists").append(
              "<option value='" + temp.id + "'>" + temp.name + "</option>");
          });
        },
        hookResponse: function (data) {
          if (!data.success) {
            $(".error").text(data.reason);
          } else {
            $("#toast").attr("text", "Webhook Created. webhook.co/" + data.handler);
            $("#toast")[0].show();
            $("#modal")[0].toggle();
            window.setTimeout(function () {
              $("iron-ajax").attr("url", "/created.json?" + data.handler);
            }, 2000);
          }
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
            method: "POST",
            url: "/trello/lists/" + value
          }).done(wh.ajax.handler.list);
        },
        createHook: function (data) {
          $.ajax({
            url: "/save",
            method: "POST",
            data: data
          }).done(wh.ajax.handler.hookResponse);
        }
      }
    },
    event: {
      handler: {
        boards: function () {
          $("#boardName").val($("option:selected", this).text());
          $(".list").show();
          var value = $(this).val();
          if (!value.length || value === "0") {
            $("#lists").html("");
            $(".list").hide();
          }
          wh.ajax.request.list(value);
        },
        lists: function () {
          $("#listName").val($("option:selected", this).text());
        },
        openCreateDialog: function () {
          $("#modal")[0].toggle();
        },
        showForm: function () {
          var service = $(this).attr("alt");
          if (service === "trello" && !$("#boards").text().trim().length) {
            wh.ajax.request.board();
          }
          $(".hide").hide();
          $("img.selected").removeClass("selected");
          $("." + service).show();
          $(".buttons .hide").show();
          if (!$(".buttons").hasClass("pad")) {
            $(".buttons").addClass("pad");
          }
          $(".error").text("");
          $(this).addClass("selected");
        },
        createHook: function () {
          var service = $("img.selected").attr("alt");
          var data = {
            "service": service
          };
          if (service === "trello") {
            data.boardName = $("#boardName").val();
            data.boardId = $("#boards").val();
            data.listName = $("#listName").val();
            data.listId = $("#lists").val();
            if (!data.boardName.length || !data.boardId.length || !data.listId.length || !data.listName.length) {
              $(".error").text("Provide all information.");
              return;
            }
          } else if (service === "telegram") {
            data.teleCode = $("#teleCode").val();
            if (data.teleCode.length !== 6) {
              $(".error").text("Invalid code.");
              return;
            }
          } else if (service === "slack") {
            data.slack_url = $("#slack").val();
            data.slack_channel = $("#slack_channel").val();
            if (data.slack_url.search("https://hooks.slack.com/services/") === -1) {
              $(".error").text("Invalid Slack URL.");
              return;
            }
            if (data.slack_channel[0] !== '@' && data.slack_channel[0] !== '#') {
              $(".error").text("Invalid Channel / Username.");
              return;
            }
          } else if (service === "pushover") {
            data.poUserkey = $("#poUserkey").val();
            if (data.poUserkey.length < 24) {
              $(".error").text("Invalid key.");
            }
          } else if (service === "hipchat") {
            data.hcToken = $("#hcToken").val();
            data.hcRoomid = $("#hcRoomid").val();
          }
          wh.ajax.request.createHook(data);
        }
      },
      add: function () {
        var handler = wh.event.handler;
        $("#lists").on("change", handler.lists);
        $("#boards").on("change", handler.boards);
        $("#addHook").on("click", handler.openCreateDialog);
        $(".app img").on("click", handler.showForm);
        $("#create").on("click", handler.createHook);
      }
    }
  };

  $(document).ready(function () {
    wh.event.add();
  });
}());
