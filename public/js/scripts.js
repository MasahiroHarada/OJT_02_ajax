$(function() {

    var currentPage = 1;

    $('.load-more').click(function () {
        currentPage += 1;
        $.ajax({
            method: 'get',
            url: '/api/gif?page=' + currentPage,
            success: function(response) {
                // アイテム追加
                var items = '';
                response.GifList.forEach(function (item) {
                    items += `<div class="gif-item">
                        <div class="content">
                            <div class="image-box">
                                <a href="/public/img/${item.Filename}" data-lightbox="${item.Filename}">
                                    <img src="/public/img/${item.Filename}">
                                </a>
                            </div>
                        </div>
                    </div>`;
                });
                $('.gif').append($(items));

                // 最後のページだった場合はもっと見るボタンを隠す
                if (response.IsLastPage) {
                    $('.load-more').hide();
                }
            }
        });
    });

    /**
     * Lightbox
     * http://lokeshdhakar.com/projects/lightbox2/
     */
    lightbox.option({
        fadeDuration: 300,
        imageFadeDuration: 300,
        resizeDuration: 400
    });

});
