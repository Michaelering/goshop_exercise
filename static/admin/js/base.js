
$(function () {
    baseApp.init();
})

var baseApp = {
    init: function () {
        this.confirmDelete()
        this.changeStatus()
        this.changeNum()
    },
    // Delete confirmation
    confirmDelete: function () {
        $(".delete").click(function () {
            var flag = confirm("您确定要删除吗？")
            return flag
        })
    },
    // Toggle status via switch
    changeStatus: function () {
        $(".chStatus").click(function () {
            var id = $(this).attr("data-id")
            var table = $(this).attr("data-table")
            var field = $(this).attr("data-field")
            var el = $(this)
            $.get("/admin/changeStatus", { id: id, table: table, field: field }, function (response) {
                if (response.success) {
                    var checkbox = el.find('input[type="checkbox"]')
                    checkbox.prop('checked', !checkbox.prop('checked'))
                }
            })
        })
    },
    // Inline number editing
    changeNum: function () {
        $(".chSpanNum").click(function () {
            var id = $(this).attr("data-id")
            var table = $(this).attr("data-table")
            var field = $(this).attr("data-field")
            var num = $(this).text().trim()
            var spanEl = $(this)

            // Already editing
            if (spanEl.find('input').length > 0) return;

            var input = $('<input type="number" style="width:70px;text-align:center;padding:2px 4px;border:2px solid #007bff;border-radius:3px;outline:none;" value="" />');
            spanEl.html(input);
            input.trigger("focus").val(num);
            input.click(function (e) {
                e.stopPropagation();
            })
            input.blur(function () {
                var inputNum = $(this).val()
                spanEl.html(inputNum)
                $.get("/admin/changeNum", { id: id, table: table, field: field, num: inputNum }, function (response) {
                    // success
                })
            })
            input.keypress(function(e) {
                if (e.which === 13) {
                    $(this).blur();
                }
            })
        })
    }
}
