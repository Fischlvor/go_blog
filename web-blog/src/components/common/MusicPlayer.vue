<template>
  <div class="music-player">
    <div class="player-background">
      <!-- 左侧旋转的唱片图标 -->
      <div class="disc" :class="{ rotating: isPlaying }">
        <img :src="currentSong.image" class="album-cover" alt="Album Cover">
        <div class="disc-center"></div>
      </div>

      <!-- 右侧控制区域 -->
      <div class="player-controls">
        <!-- 上部歌名显示 -->
        <div class="song-info">
          <div class="song-info-scroll">
            <span class="song-name">{{ currentSong.name }}</span>
            <span class="separator"> - </span>
            <span class="artist-name">{{ currentSong.artist }}</span>
          </div>
        </div>

        <!-- 下部控制按钮 -->
        <!-- 下部控制按钮 - 添加了右移的容器 -->
        <div class="controls-container">
          <div class="controls">
            <button class="control-btn" @click="playPrevious" title="上一首">
              <svg viewBox="0 0 24 24" width="16" height="16">
                <path d="M6 6h2v12H6zm3.5 6l8.5 6V6z" fill="currentColor"/>
              </svg>
            </button>
            <button class="control-btn play-btn" @click="togglePlay" :title="isPlaying ? '暂停' : '播放'">
              <svg viewBox="0 0 24 24" width="18" height="18">
                <path v-if="isPlaying" d="M6 19h4V5H6v14zm8-14v14h4V5h-4z" fill="currentColor"/>
                <path v-else d="M8 5v14l11-7z" fill="currentColor"/>
              </svg>
            </button>
            <button class="control-btn" @click="playNext" title="下一首">
              <svg viewBox="0 0 24 24" width="16" height="16">
                <path d="M6 18l8.5-6L6 6v12zM16 6v12h2V6h-2z" fill="currentColor"/>
              </svg>
            </button>

            <!-- 音量控制 -->
            <div class="volume-control">
              <button class="volume-btn" title="音量">
                <svg viewBox="0 0 24 24" width="14" height="14">
                  <path d="M3 9v6h4l5 5V4L7 9H3zm13.5 3c0-1.77-1.02-3.29-2.5-4.03v8.05c1.48-.73 2.5-2.25 2.5-4.02zM14 3.23v2.06c2.89.86 5 3.54 5 6.71s-2.11 5.85-5 6.71v2.06c4.01-.91 7-4.49 7-8.77s-2.99-7.86-7-8.77z" fill="currentColor"/>
                </svg>
              </button>
              <div class="volume-slider">
                <input type="range"
                       min="0"
                       max="1"
                       step="0.01"
                       v-model="volume"
                       @input="handleVolumeChange"
                       class="slider">
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <audio ref="audioPlayer"
           @play="isPlaying = true"
           @pause="isPlaying = false"
           @ended="playNext"
           preload="none"></audio>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, nextTick } from 'vue';

interface Song {
  name: string;
  artist: string;
  url: string;
  image: string;
}

const songs = ref<Song[]>([
  {
    name: 'Moon Halo',
    artist: '茶理理/TetraCalyx/Hanser/HOYO-MiX',
    url: 'https://music.163.com/song/media/outer/url?id=1859652717.mp3',
    image: 'http://p2.music.126.net/ciLKATqflV2YWSV3ltE2Kw==/109951166159281275.jpg'
  },
]);

const currentSongIndex = ref(0);
const audioPlayer = ref<HTMLAudioElement | null>(null);
const isPlaying = ref(false);
const volume = ref(0.1);

const currentSong = computed(() => songs.value[currentSongIndex.value]);
const currentSongUrl = computed(() => currentSong.value?.url || '');

onMounted(() => {
  if (audioPlayer.value) {
    audioPlayer.value.volume = volume.value;
  }
});

function handleVolumeChange() {
  if (audioPlayer.value) {
    audioPlayer.value.volume = volume.value;
  }
}

function togglePlay() {
  if (!audioPlayer.value) return;

  if (isPlaying.value) {
    audioPlayer.value.pause();
  } else {
    // 懒加载：只有在播放时才设置src并下载歌曲
    if (!audioPlayer.value.src || audioPlayer.value.src !== currentSongUrl.value) {
      audioPlayer.value.src = currentSongUrl.value;
    }
    audioPlayer.value.play().catch(error => {
      console.error('播放失败:', error);
    });
  }
}

function playPrevious() {
  let newIndex = currentSongIndex.value - 1;
  if (newIndex < 0) newIndex = songs.value.length - 1;
  playSong(newIndex);
}

function playNext() {
  let newIndex = currentSongIndex.value + 1;
  if (newIndex >= songs.value.length) newIndex = 0;
  playSong(newIndex);
}

function playSong(index: number) {
  currentSongIndex.value = index;

  nextTick(() => {
    if (audioPlayer.value) {
      // 懒加载：切换歌曲时才设置新的src
      audioPlayer.value.src = currentSongUrl.value;
      audioPlayer.value.currentTime = 0;
      audioPlayer.value.play().catch(error => {
        console.error('播放失败:', error);
      });
    }
  });
}
</script>

<style scoped lang="scss">
/* 音乐播放器最外层容器样式 */
.music-player {
  font-family: 'Arial', sans-serif;  /* 设置字体 */
  display: flex;                     /* 启用flex布局 */
  justify-content: center;           /* 水平居中 */
  align-items: center;               /* 垂直居中 */

  /* 播放器主体区域（包含唱片和控制区） */
  .player-background {
    display: flex;                     /* 启用flex布局 */
    align-items: center;               /* 垂直居中 */
    background: #333333;               /* 浅灰色背景 */
    border-radius: 16px;               /* 圆角边框 */
    padding: 8px;                     /* 内边距 */
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1); /* 阴影效果 */
    width: 320px;                      /* 固定宽度 */
    height: 88px;

    /* 唱片样式 */
    .disc {
      width: 72px;                       /* 唱片宽度 */
      height: 72px;                      /* 唱片高度 */
      border-radius: 50%;                /* 圆形 */
      position: relative;                /* 相对定位 */
      //margin-left: 0px;                /* 右侧外边距 */
      margin-right: 12px;                /* 右侧外边距 */
      display: flex;
      justify-content: center;           /* 水平居中 */
      align-items: center;               /* 垂直居中 */
      box-shadow: 0 2px 6px rgba(0, 0, 0, 0.2); /* 阴影 */
      overflow: hidden;                  /* 隐藏溢出内容 */

      /* 播放时的旋转动画类 */
      &.rotating {
        animation: rotate 10s linear infinite; /* 10秒无限循环旋转 */
      }

      /* 唱片封面图片样式 */
      .album-cover {
        width: 100%;                       /* 宽度100% */
        height: 100%;                      /* 高度100% */
        object-fit: cover;                 /* 图片填充方式 */
        border-radius: 50%;                /* 圆形 */
      }

      /* 唱片中心圆点样式 */
      .disc-center {
        width: 16px;                       /* 直径 */
        height: 16px;                      /* 直径 */
        border-radius: 50%;                /* 圆形 */
        background: white;                 /* 白色背景 */
        position: absolute;                /* 绝对定位 */
        z-index: 2;                        /* 层级 */
      }

    }

    /* 右侧控制区域容器 */
    .player-controls {
      flex: 1;                           /* 占据剩余空间 */
      display: flex;
      flex-direction: column;            /* 垂直排列 */
      justify-content: space-between;    /* 两端对齐 */
      //height: 80px;                       /* 固定高度 */
      //gap: 1px; /* 新增：控制上下元素间距，可根据需要调整 */
      height: auto; /* 改为自动高度 */

      /* 歌曲信息区域 */
      .song-info {
        text-align: left;                  /* 左对齐 */
        margin-bottom: 8px;                 /* 下边距 */
        overflow: hidden;                  /* 隐藏溢出 */
        white-space: nowrap;               /* 不换行 */
        width: 100%;                       /* 宽度100% */
        position: relative;                /* 相对定位 */
        height: 24px;                      /* 固定高度 */

        /* 滚动文本容器 */
        .song-info-scroll {
          display: inline-block;           /* 行内块元素 */
          position: absolute;               /* 绝对定位 */
          white-space: nowrap;             /* 不换行 */
          padding-left: 100%;              /* 初始位置在右侧 */
          animation: scrollText 15s linear infinite; /* 15秒滚动动画 */

          /* 鼠标悬停暂停动画 */
          &:hover {
            animation-play-state: paused;
          }
        }

        /* 歌曲名称样式 */
        .song-name {
          font-size: 16px;                 /* 字体大小 */
          font-weight: bold;               /* 加粗 */
          color: #DCDFE6;                     /* 深灰色 */
        }

        /* 艺术家名称样式 */
        .artist-name {
          font-size: 14px;                 /* 字体大小 */
          color:#DCDFE6;                     /* 中灰色 */
        }

        /* 分隔符样式 */
        .separator {
          margin: 0 4px;                   /* 左右边距 */
          color: #DCDFE6;                     /* 中灰色 */
        }
      }
    }

    /* 控制按钮容器（新增） */
    .controls-container {
      display: flex;
      justify-content: center;         /* 内容靠右对齐 */
      //padding-right: 20px;               /* 右侧内边距（控制右移距离） */

      /* 控制按钮组样式 */
      .controls {
        display: flex;
        align-items: center;               /* 垂直居中 */
        gap: 8px;                          /* 按钮间距 */

        /* 控制按钮基础样式 */
        .control-btn {
          background: none;                  /* 透明背景 */
          border: none;                      /* 无边框 */
          width: 32px;                       /* 宽度 */
          height: 32px;                      /* 高度 */
          border-radius: 50%;                /* 圆形 */
          display: flex;
          justify-content: center;           /* 水平居中 */
          align-items: center;               /* 垂直居中 */
          cursor: pointer;                   /* 手型光标 */
          color: white;                       /* 按钮颜色 */
          transition: all 0.2s;              /* 过渡动画 */
          padding: 0;                        /* 无内边距 */

          /* 悬停效果 */
          &:hover {
            background: #e9ecef;             /* 浅灰色背景 */
            color: #333;                     /* 深灰色文字 */
          }

          /* 播放按钮特殊样式 */
          &.play-btn {
            background: #4CAF50;             /* 绿色背景 */
            color: white;                    /* 白色文字 */

            /* 播放按钮悬停效果 */
            &:hover {
              background: #3e8e41;           /* 深绿色 */
            }
          }

          /* 图标大小同步调整 */
          svg {
            width: 25px;   /* 标准按钮图标大小 */
            height: 25px;
          }

          /* 播放/暂停按钮特殊样式（绿色大按钮） */
          &.play-btn {
            width: 36px;   /* 比标准按钮大约12% */
            height: 36px;
          }

          /* 播放按钮图标稍大 */
          &.play-btn svg {
            width: 25px;
            height: 25px;
          }
        }

        /* 音量控制区域 */
        .volume-control {
          position: relative;                /* 相对定位 */
          display: flex;
          align-items: center;               /* 垂直居中 */

          /* 音量按钮样式 */
          .volume-btn {
            background: none;                /* 透明背景 */
            border: none;                    /* 无边框 */
            cursor: pointer;                 /* 手型光标 */
            color: white;                     /* 按钮颜色 */
            display: flex;
            align-items: center;             /* 垂直居中 */
            justify-content: center;         /* 水平居中 */
            width: 24px;                     /* 宽度 */
            height: 24px;                     /* 高度 */
            padding: 0;                      /* 无内边距 */

            /* 音量滑块显示逻辑 */
            &:hover + .volume-slider,
            .volume-slider:hover {
              opacity: 1;                    /* 完全显示 */
              width: 150px;                   /* 展开宽度 */
            }
          }

          /* 播放按钮图标稍大 */
          .volume-btn svg {
            width: 25px;
            height: 25px;
          }

          /* 音量滑块容器 */
          .volume-slider {
            position: absolute;              /* 绝对定位 */
            left: -80px;                      /* 距离左侧 */
            top: 25px;
            opacity: 0;                      /* 初始透明 */
            width: 0;                        /* 初始宽度为0 */
            transition: all 0.3s ease;       /* 过渡动画 */
            background: white;               /* 白色背景 */
            padding: 4px 8px;                /* 内边距 */
            border-radius: 16px;             /* 圆角 */
            box-shadow: 0 1px 3px rgba(0,0,0,0.1); /* 阴影 */

            /* 滑块输入框样式 */
            .slider {
              width: 100%;                   /* 宽度100% */
              cursor: pointer;               /* 手型光标 */
              height: 4px;                   /* 高度 */
            }
          }

          /* 音量控制区域悬停效果 */
          &:hover .volume-slider {
            opacity: 1;                      /* 显示滑块 */
            width: 150px;                     /* 展开宽度 */
          }
        }
      }
    }

  }
}











/* 唱片旋转动画定义 */
@keyframes rotate {
  from { transform: rotate(0deg); }  /* 起始角度 */
  to { transform: rotate(360deg); }  /* 结束角度 */
}

/* 文本滚动动画定义 */
@keyframes scrollText {
  0% {
    transform: translateX(0);        /* 起始位置 */
  }
  100% {
    transform: translateX(-100%);   /* 向左移动100% */
  }
}











/* 隐藏audio元素 */
audio {
  display: none;
}
</style>