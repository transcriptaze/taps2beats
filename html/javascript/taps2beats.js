'use strict';

var player
var loopTimer
var delayTimer
var loaded = false
var looping = false
var start
var end
var taps = {
  duration: 0,
  taps: [],
  current: []
}

function onPlayerReady(event) {
  document.getElementById('url').readOnly = false
  document.getElementById('load').disabled = false

  start = new Slider('start', onSetStart)
  end = new Slider('end', onSetEnd)
        end.init(0,100,100)
}

function onPlayerStateChange(event) {
  clearInterval(loopTimer)

  switch (event.data) {
    case YT.PlayerState.ENDED:
      if (!loaded) {
        loaded = true
        const duration = player.getDuration()

        start.init(0,duration,0)
        end.init(0,duration,duration)
        taps.duration = duration
        cue(false)
      }
      break

    case YT.PlayerState.PLAYING:
      if (loaded) {
        loopTimer = setInterval(tick, 100)        
        document.getElementById('tap').dataset.state = 'playing'
      }
      break

    case YT.PlayerState.PAUSED:
      document.getElementById('tap').dataset.state = 'cued'
      break

    case YT.PlayerState.BUFFERING:
      break

    case YT.PlayerState.CUED:
      document.getElementById('loading').style.visibility = 'hidden'
      document.getElementById('controls').style.visibility = 'visible'      
      document.getElementById('tap').style.visibility = 'visible'      
      document.getElementById('taps').style.visibility = 'visible'      
      document.getElementById('tap').dataset.state = 'cued'
      document.getElementById('pad').focus()
      player.unMute()

      if (taps.current.length > 0)  {
          taps.taps.push(taps.current)
          taps.current = []
          draw()
      }
      break
  }
}

function onSetStart(t, released) {
  document.getElementById('from').value = format(t)

  switch (player.getPlayerState()) {
    case YT.PlayerState.CUED:
    case YT.PlayerState.ENDED:
      if (released) {
        cue(false)
      }
      break

    case YT.PlayerState.PAUSED:
      player.seekTo(t, released)                              
      break;

    default:
      if (released && player.getCurrentTime() < t) {
          player.seekTo(t, true)                              
      }
    }
}

function onSetEnd(t, released) {
  document.getElementById('to').value = format(t)

  if (released) {
    switch (player.getPlayerState()) {
      case YT.PlayerState.ENDED:
      case YT.PlayerState.CUED:
        cue(false)
        break
    }    
  }
}

function onLoop(event) {
  looping = event.target.checked
}

function onTap(event) {
  if (event.code === 'Space') {
    event.preventDefault()

    if (!event.repeat) {
      switch (player.getPlayerState()) {
        case YT.PlayerState.CUED:
        case YT.PlayerState.PAUSED:
          player.playVideo()
          break

        case YT.PlayerState.PLAYING:
          taps.current.push(player.getCurrentTime())
          draw()
          break
      }
    }
  } else if (event.code === 'KeyS') {
    event.preventDefault()
    if (!event.repeat && player.getPlayerState() == YT.PlayerState.PLAYING) {
        cue(false)
    }
  } else if (event.code === 'KeyP') {
    event.preventDefault()
    if (!event.repeat && player.getPlayerState() == YT.PlayerState.PLAYING) {
        player.pauseVideo()
    }
  } 
}

function onFocus(event) {
  switch (player.getPlayerState()) {
    case YT.PlayerState.CUED:
    case YT.PlayerState.PAUSED:
      document.getElementById('tap').dataset.state = 'cued'
      break

    case YT.PlayerState.PLAYING:
      document.getElementById('tap').dataset.state = 'playing'
  }
}

function onBlur(event) {
  document.getElementById('tap').dataset.state = 'idle'
}

function load(event) {
  const url = document.getElementById('url')    
  const vid = getVideoID(url.value)

  loaded = false
  document.getElementById('loading').style.visibility = 'visible'
  player.loadVideoById({ videoId: vid, startSeconds: 0, endSeconds: 0.1 })
}

function analyse (event) {
  const p = new Promise((resolve, reject) => {
    goTaps((obj, err) => {
      if (err) {
        reject(err)
      } else {
        resolve(obj)
      }
    })
  })

  p.then(o => { console.log("gotcha!!", o) })
   .catch(function (err) { console.error(err) })
}

function cue(play) {
  const url = document.getElementById('url')    
  const vid = getVideoID(url.value)
  const start = document.getElementById('start').getAttribute('aria-valuenow')

  if (play) {
    player.loadVideoById({ videoId: vid, startSeconds: start })
  } else {
    player.mute()
    player.cueVideoById({ videoId: vid, startSeconds: start })
  }
}

function tick() {
  const delay = parseFloat(document.getElementById('delay').value)
  const end = document.getElementById('end').getAttribute('aria-valuenow')
  const t = player.getCurrentTime()

  if (t > end) {
    if (!isNaN(delay) && delay > 0) {
        cue(false)      
        delayTimer = setInterval(delayed, delay * 1000)
    } else {
        cue(looping)            
    }
  }
}

function delayed() {
  clearInterval(delayTimer)
  cue(looping)
}

function draw() {
  if (player.getPlayerState() ===  YT.PlayerState.PLAYING) {
      drawTaps(document.querySelector('#current canvas.all'), taps.current, 0, taps.duration)
      drawTaps(document.querySelector('#current canvas.zoomed'), taps.current, start.valueNow, end.valueNow - start.valueNow)
  } else {
      drawTaps(document.querySelector('#history canvas.all'), taps.taps, 0, taps.duration)    
      drawTaps(document.querySelector('#history canvas.zoomed'), taps.taps, start.valueNow, end.valueNow - start.valueNow)
  }
}

function drawTaps(canvas, taps, offset, duration) {
  let ctx = canvas.getContext('2d')
  let width = canvas.width
  let height = canvas.height

  ctx.lineWidth = 1
  ctx.strokeStyle = 'red'

  ctx.clearRect(0, 0, width, height)
  ctx.beginPath()
  taps.flat().forEach(function(t) {
      let x = Math.floor((t - offset) * width/duration) + 0.5
      ctx.moveTo(x, 0)
      ctx.lineTo(x, height)
  })
  ctx.stroke()
}

var Slider = function (node, handler) {
  this.domNode = document.getElementById(node)
  this.railDomNode = document.getElementById(node).parentNode
  this.handler = handler

  this.minDomNode = false
  this.maxDomNode = false

  this.valueNow = 50

  this.railMin = 0
  this.railMax = 100
  this.railWidth = 0

  this.thumbWidth = 10

  this.keyCode = Object.freeze({
    left: 37,
    up: 38,
    right: 39,
    down: 40,
    pageUp: 33,
    pageDown: 34,
    end: 35,
    home: 36,
  })
}

// Initialize slider
Slider.prototype.init = function (min, max, t) {
  this.valueNow = t
  this.domNode.setAttribute('aria-valuenow', t)
  this.domNode.setAttribute('aria-valuemin', min)
  this.domNode.setAttribute('aria-valuemax', max)

  this.railMin = min
  this.railMax = max
  this.railWidth = parseInt(this.railDomNode.style.width.slice(0, -2));

  if (this.domNode.previousElementSibling) {
    this.minDomNode = this.domNode.previousElementSibling;
  }

  if (this.domNode.nextElementSibling) {
    this.maxDomNode = this.domNode.nextElementSibling;
  }

  if (this.domNode.tabIndex != 0) {
    this.domNode.tabIndex = 0;
  }

  this.domNode.addEventListener('keydown', this.handleKeyDown.bind(this));
  this.domNode.addEventListener('mousedown', this.handleMouseDown.bind(this));
  this.domNode.addEventListener('focus', this.handleFocus.bind(this));
  this.domNode.addEventListener('blur', this.handleBlur.bind(this));

  this.moveSliderTo(this.valueNow);
};

Slider.prototype.moveSliderTo = function (value, released) {
  var valueMax = parseInt(this.domNode.getAttribute('aria-valuemax'));
  var valueMin = parseInt(this.domNode.getAttribute('aria-valuemin'));

  if (value > valueMax) {
    value = valueMax;
  }

  if (value < valueMin) {
    value = valueMin
  }

  this.valueNow = value
  this.dolValueNow = value

  if (released) {
    this.domNode.setAttribute('aria-valuenow', this.valueNow);    
  }

  if (this.minDomNode) {
    this.minDomNode.setAttribute('aria-valuemax', this.valueNow);
  }

  if (this.maxDomNode) {
    this.maxDomNode.setAttribute('aria-valuemin', this.valueNow);
  }

  const range = this.railMax - this.railMin
  const scale = this.railWidth/range
  const pos = (this.valueNow - this.railMin) * scale

  if (this.minDomNode) {
    this.domNode.style.left = Math.ceil(pos) + 'px'
  } else {
    this.domNode.style.left = Math.floor(pos) - this.thumbWidth + 'px'
  }

  if (this.handler) {
    this.handler(this.dolValueNow, released)
  }
}

Slider.prototype.handleKeyDown = function (event) {
  var flag = false;

  switch (event.keyCode) {
    case this.keyCode.left:
    case this.keyCode.down:
      this.moveSliderTo(this.valueNow - 1, true)
      flag = true;
      break;

    case this.keyCode.right:
    case this.keyCode.up:
      this.moveSliderTo(this.valueNow + 1, true)
      flag = true;
      break;

    case this.keyCode.pageDown:
      this.moveSliderTo(this.valueNow - 10, true)
      flag = true;
      break;

    case this.keyCode.pageUp:
      this.moveSliderTo(this.valueNow + 10, true)
      flag = true;
      break;

    case this.keyCode.home:
      this.moveSliderTo(this.railMin, true)
      flag = true;
      break;

    case this.keyCode.end:
      this.moveSliderTo(this.railMax, true)
      flag = true;
      break;

    default:
      break;
  }

  if (flag) {
    event.preventDefault();
    event.stopPropagation();
  }
};

Slider.prototype.handleFocus = function () {
  this.domNode.classList.add('focus');
  this.railDomNode.classList.add('focus');
};

Slider.prototype.handleBlur = function () {
  this.domNode.classList.remove('focus');
  this.railDomNode.classList.remove('focus');
};

Slider.prototype.handleMouseDown = function (event) {
  var self = this;

  var handleMouseMove = function (event) {
    var diffX = event.pageX - self.railDomNode.offsetLeft;
    self.valueNow =
      self.railMin +
      parseInt(((self.railMax - self.railMin) * diffX) / self.railWidth);

    self.moveSliderTo(self.valueNow);

    event.preventDefault();
    event.stopPropagation();
  };

  var handleMouseUp = function () {
    document.removeEventListener('mousemove', handleMouseMove);
    document.removeEventListener('mouseup', handleMouseUp);

    self.moveSliderTo(self.valueNow, true);
  };

  // bind a mousemove event handler to move pointer
  document.addEventListener('mousemove', handleMouseMove);

  // bind a mouseup event handler to stop tracking mouse movements
  document.addEventListener('mouseup', handleMouseUp);

  event.preventDefault();
  event.stopPropagation();

  // Set focus to the clicked handle
  this.domNode.focus();
};

function format(t) {
  let minutes = 0
  let seconds = 0

  if (t > 0) {
    minutes = Math.floor(t/60)
    seconds = t % 60
  }

  return String(minutes) + ':' + String(seconds).padStart(2,'0')
}

function getVideoID(url) {
  let addr = new URL(url)

  try {
    return addr.searchParams.get("v")
  } catch(err) {
    console.log(err)
  }

  return ""
}

