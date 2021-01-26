'use strict';

var player
var loopTimer
var cued = false
var start
var end

function onPlayerReady(event) {
  document.getElementById('url').readOnly = false
  document.getElementById('load').disabled = false

  start = new Slider('start', onSetStart)
  end = new Slider('end', onSetEnd)
}

function onPlayerStateChange(event) {
  clearInterval(loopTimer)

  switch (event.data) {
    case YT.PlayerState.ENDED:
      if (!cued) {
        cued = true
        const duration = player.getDuration()

        start.init(0,duration,0)
        end.init(0,duration,duration)
        cue()
      }
      break

    case YT.PlayerState.PLAYING:
      if (cued) {
        loopTimer = setInterval(tick, 100)        
      }
      break

    case YT.PlayerState.PAUSED:
      break

    case YT.PlayerState.BUFFERING:
      break

    case YT.PlayerState.CUED:
      document.getElementById('loading').style.visibility = 'hidden'
      document.getElementById('controls').style.visibility = 'visible'      
      break
  }
}

function tick() {
  const end = document.getElementById('end').getAttribute('aria-valuenow')
  const t = player.getCurrentTime()

  if (t > end) {
    cue()
  }
}

function onSetStart(t, released) {
  document.getElementById('from').value = format(t)

  switch (player.getPlayerState()) {
    case YT.PlayerState.CUED:
    case YT.PlayerState.ENDED:
      if (released) {
        cue()
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
        cue()
        break
    }    
  }
}

function load(event) {
  const url = document.getElementById('url')    
  const vid = getVideoID(url.value)

  cued = false
  document.getElementById('loading').style.visibility = 'visible'
  player.loadVideoById({ videoId: vid, startSeconds: 0, endSeconds: 0.1 })
}

function cue() {
  const url = document.getElementById('url')    
  const vid = getVideoID(url.value)
  const start = document.getElementById('start').getAttribute('aria-valuenow')

  player.cueVideoById({ videoId: vid, startSeconds: start })
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
  this.railBorderWidth = 1

  this.thumbWidth = 20
  this.thumbHeight = 24

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
    value = valueMin;
  }

  this.valueNow = value;
  this.dolValueNow = value;

  if (released) {
    this.domNode.setAttribute('aria-valuenow', this.valueNow);    
  }

  if (this.minDomNode) {
    this.minDomNode.setAttribute('aria-valuemax', this.valueNow);
  }

  if (this.maxDomNode) {
    this.maxDomNode.setAttribute('aria-valuemin', this.valueNow);
  }

  var pos = Math.round(
    ((this.valueNow - this.railMin) *
      (this.railWidth - 2 * (this.thumbWidth - this.railBorderWidth))) /
      (this.railMax - this.railMin)
  );

  if (this.minDomNode) {
    this.domNode.style.left = pos + this.thumbWidth - this.railBorderWidth + 'px';
  } else {
    this.domNode.style.left = pos - this.railBorderWidth + 'px';
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

