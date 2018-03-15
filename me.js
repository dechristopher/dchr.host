console.log(`© 2018 Andrew DeChristopher ; contact for inquiries or job opportunities`);
document.addEventListener('touchstart', onTouchStart, { passive: true });

$('.to').textillate({ in: { effect: 'fadeInDown', shuffle: true } });

setTimeout(function () {
	document.getElementById('me').innerHTML = document.getElementById('me').innerHTML + '\n<i class="y">why are you still here?</i>';
	$('.y').textillate({ in: { effect: 'fadeInUp' } });
}, 35000);

