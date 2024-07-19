document.addEventListener('scroll', () => {
    const boxes = document.querySelectorAll('.box');
    const windowHeight = window.innerHeight;

    boxes.forEach((box, index) => {
        let boxTop = box.getBoundingClientRect().top;
        let video = box.querySelector('video');

        if (boxTop < windowHeight && boxTop > -windowHeight) {
            // Current video is in view
            video.style.transform = 'scale(1)';
            video.play();
        } else if (boxTop <= -windowHeight) {
            // Previous video
            video.style.transform = 'scale(0.8)';
            video.pause();
        } else {
            // Upcoming video
            video.style.transform = 'scale(0.8)';
            video.pause();
        }
    });
});
