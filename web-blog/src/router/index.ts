import { createRouter, createWebHistory } from 'vue-router'
import {useUserStore} from "@/stores/user";
import {useLayoutStore} from "@/stores/layout";
import { ElMessage, ElMessageBox } from 'element-plus';

const routes = [
  {
    path: '/',
    name: 'web',
    component: () => import('@/views/web/index.vue'),
    children: [
      {
        path: "/",
        name: "index",
        component: () => import('@/views/web/index/index.vue'),
        meta: {
          title: "首页"
        }
      },
      {
        path: "search",
        name: "search",
        component: () => import('@/views/web/search/index.vue'),
        meta: {
          title: "搜索"
        }
      },
      {
        path: "news",
        name: "news",
        component: () => import('@/views/web/news/index.vue'),
        meta: {
          title: "新闻"
        }
      },
      {
        path: "friend-link",
        name: "friend-link",
        component: () => import('@/views/web/friend-link/index.vue'),
        meta: {
          title: "友链"
        }
      },
      {
        path: "about",
        name: "about",
        component: () => import('@/views/web/about/index.vue'),
        meta: {
          title: "关于"
        }
      },
      {
        path: "/article/:id",
        name: "article",
        component: () => import('@/views/web/article/index.vue'),
        meta: {
          title: "文章"
        }
      },
      {
        path: "ai-assistant",
        name: "ai-assistant",
        component: () => import('@/views/web/ai-assistant/index.vue'),
        meta: {
          title: "AI助手",
          requiresAuth: true
        }
      }

    ]
  },
  {
    path: "/login",
    name: "login",
    component: () => import('@/views/login/index.vue')
  },
  {
    path: "/sso-callback",
    name: "sso-callback",
    component: () => import('@/views/SSOCallback.vue'),
    meta: {
      title: "登录处理中"
    }
  },
  // {
  //   path: "/article/:id",
  //   name: "article",
  //   component: () => import('@/views/web/article/index.vue')
  // },
  {
    path: "/dashboard",
    name: "dashboard",
    component: () => import('@/views/dashboard/index.vue'),
    meta: {
      title: "控制面板",
      requiresAuth: true
    },
    children: [
      {
        path: "/dashboard/",
        name: "home",
        component: () => import('@/views/dashboard/home/index.vue'),
        meta: {
          title: "主页"
        },
      },
      {
        path: "user-center",
        name: "user-center",
        meta: {
          title: "个人中心"
        },
        children: [
          {
            path: "user-info",
            name: "user-info",
            component: () => import('@/views/dashboard/user-center/user-info.vue'),
            meta: {
              title: "我的信息"
            }

          },
          {
            path: "user-star",
            name: "user-star",
            component: () => import('@/views/dashboard/user-center/user-star.vue'),
            meta: {
              title: "我的收藏"
            }
          },
          {
            path: "user-comment",
            name: "user-comment",
            component: () => import('@/views/dashboard/user-center/user-comment.vue'),
            meta: {
              title: "我的评论"
            }
          },
          {
            path: "user-feedback",
            name: "user-feedback",
            component: () => import('@/views/dashboard/user-center/user-feedback.vue'),
            meta: {
              title: "我的反馈"
            }
          }
        ]
      },
      {
        path: "users",
        name: "users",
        meta: {
          title: "用户管理",
          requiresAdmin: true
        },
        children: [
          {
            path: "user-list",
            name: "user-list",
            component: () => import('@/views/dashboard/users/user-list.vue'),
            meta: {
              title: "用户列表"
            }
          }
        ]
      },
      {
        path: "articles",
        name: "articles",
        meta: {
          title: "文章管理",
          requiresAdmin: true
        },
        children: [
          {
            path: "article-publish",
            name: "article-publish",
            component: () => import('@/views/dashboard/articles/article-publish.vue'),
            meta: {
              title: "发布文章"
            }
          },
          {
            path: "comment-list",
            name: "comment-list",
            component: () => import('@/views/dashboard/articles/comment-list.vue'),
            meta: {
              title: "评论列表"
            }
          },
          {
            path: "article-list",
            name: "article-list",
            component: () => import('@/views/dashboard/articles/article-list.vue'),
            meta: {
              title: "文章列表"
            }
          }
        ]
      },
      {
        path: "images",
        name: "images",
        meta: {
          title: "图片管理",
          requiresAdmin: true
        },
        children: [
          {
            path: "image-list",
            name: "image-list",
            component: () => import('@/views/dashboard/images/image-list.vue'),
            meta: {
              title: "图片列表"
            }
          }
        ]
      },
      {
        path: "emoji",
        name: "emoji",
        meta: {
          title: "表情管理",
          requiresAdmin: true
        },
        children: [
          {
            path: "emoji-list",
            name: "emoji-list",
            component: () => import('@/views/dashboard/emoji/emoji-list.vue'),
            meta: {
              title: "表情列表"
            }
          },
          {
            path: "emoji-groups",
            name: "emoji-groups",
            component: () => import('@/views/dashboard/emoji/emoji-groups.vue'),
            meta: {
              title: "表情组管理"
            }
          },
          {
            path: "emoji-sprites",
            name: "emoji-sprites",
            component: () => import('@/views/dashboard/emoji/emoji-sprites.vue'),
            meta: {
              title: "雪碧图管理"
            }
          }
        ]
      },
      {
        path: "system",
        name: "system",
        meta: {
          title: "系统管理",
          requiresAdmin: true
        },
        children: [
          {
            path: "feedback-list",
            name: "feedback-list",
            component: () => import('@/views/dashboard/system/feedback-list.vue'),
            meta: {
              title: "反馈列表"
            }
          },
          {
            path: "advertisement-list",
            name: "advertisement-list",
            component: () => import('@/views/dashboard/system/advertisement-list.vue'),
            meta: {
              title: "广告列表"
            }
          },
          {
            path: "friend-link-list",
            name: "friend-link-list",
            component: () => import('@/views/dashboard/system/friend-link-list.vue'),
            meta: {
              title: "友链列表"
            }
          },
          {
            path: "login-logs",
            name: "login-logs",
            component: () => import('@/views/dashboard/system/login-logs.vue'),
            meta: {
              title: "登录日志"
            }
          },
          {
            path: "app-config",
            name: "app-config",
            redirect: "/dashboard/system/app-config/site-config",
            component: () => import('@/views/dashboard/system/app-config.vue'),
            meta: {
              title: "应用配置"
            },
            children: [
              {
                path: "site-config",
                name: "site-config",
                component: () => import('@/views/dashboard/system/config/site-config.vue'),
                meta: {
                  title: "网站配置"
                }
              },
              {
                path: "system-config",
                name: "system-config",
                component: () => import('@/views/dashboard/system/config/system-config.vue'),
                meta: {
                  title: "系统配置"
                }
              },
              {
                path: "email-config",
                name: "email-config",
                component: () => import('@/views/dashboard/system/config/email-config.vue'),
                meta: {
                  title: "邮箱配置"
                }
              },
              {
                path: "qq-config",
                name: "qq-config",
                component: () => import('@/views/dashboard/system/config/qq-config.vue'),
                meta: {
                  title: "QQ登录配置"
                }
              },
              {
                path: "qiniu-config",
                name: "qiniu-config",
                component: () => import('@/views/dashboard/system/config/qiniu-config.vue'),
                meta: {
                  title: "七牛云配置"
                }
              },
              {
                path: "jwt-config",
                name: "jwt-config",
                component: () => import('@/views/dashboard/system/config/jwt-config.vue'),
                meta: {
                  title: "JWT配置"
                }
              },
              {
                path: "gaode-config",
                name: "gaode-config",
                component: () => import('@/views/dashboard/system/config/gaode-config.vue'),
                meta: {
                  title: "高德API配置"
                }
              }
            ]
          }
        ]
      },
      {
        path: "ai-management",
        name: "ai-management",
        meta: {
          title: "AI对话管理",
          requiresAdmin: true
        },
        children: [
          {
            path: "models",
            name: "ai-models",
            component: () => import('@/views/dashboard/ai-management/models.vue'),
            meta: {
              title: "模型管理"
            }
          },
          {
            path: "sessions",
            name: "ai-sessions",
            component: () => import('@/views/dashboard/ai-management/sessions.vue'),
            meta: {
              title: "会话管理"
            }
          },
          {
            path: "messages",
            name: "ai-messages",
            component: () => import('@/views/dashboard/ai-management/messages.vue'),
            meta: {
              title: "消息管理"
            }
          }
        ]
      }
    ]
  },
  {
    path: "/404",
    name: "404",
    component: () => import('@/views/error/index.vue')
  },
  {
    path: "/:catchAll(.*)",
    component: () => import('@/views/error/index.vue')
  }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: routes,
})

export default router

router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  const layoutStore = useLayoutStore()
  userStore.initializeUserInfo().then(() => {
    const isAuthenticated = userStore.isLoggedIn // 检查用户是否登录的逻辑
    const isAdmin = userStore.isAdmin // 检查用户是否为管理员的逻辑
    if (to.matched.some(record => record.meta.requiresAuth)) {
      if (!isAuthenticated) {
        // 根据目标页面显示不同的提示信息
        const pageName = to.meta.title || '该页面'
        const message = to.name === 'ai-assistant' 
          ? '✨ AI助手需要登录后才能使用哦！登录后即可开启智能对话体验～'
          : `${pageName}需要登录后才能访问，是否前往登录？`
        
        ElMessageBox.confirm(message, '需要登录', {
              cancelButtonText: '取消',
              confirmButtonText: '立即登录',
              type: 'info',
              center: true
            })
            .then(async () => {
              // 直接跳转到SSO登录页面
              try {
                const redirectUri = encodeURIComponent(window.location.origin + '/sso-callback');
                const response = await fetch(`/api/auth/sso_login_url?redirect_uri=${redirectUri}`);
                const data = await response.json();
                
                if (data.code === 0) {
                  // 使用state作为key，存储返回URL到sessionStorage
                  const state = data.data.state;
                  sessionStorage.setItem(`oauth_state_${state}`, JSON.stringify({
                    returnUrl: window.location.pathname,
                    timestamp: Date.now()
                  }));
                  
                  window.location.href = data.data.sso_login_url;
                } else {
                  ElMessage.error(data.message || '获取登录地址失败');
                  router.push({name: 'index', replace: true});
                }
              } catch (error) {
                console.error('获取SSO登录URL失败:', error);
                ElMessage.error('登录服务异常，请稍后重试');
                router.push({name: 'index', replace: true});
              }
            })
            .catch(() => {
              router.push({name: from.name as string}).then();
            });
      } else if (to.matched.some(record => record.meta.requiresAdmin) && !isAdmin) {
        ElMessageBox.confirm(
            '权限不足，请确认您的用户角色是否具备访问该页面的权限。', 'Warning', {
              confirmButtonText: '确定',
              type: 'warning',
            })
            .then(() => {
              router.push({name: from.name as string}).then();
            });
      } else {
        next(); // 继续访问
      }
    } else {
      next(); // 不需要登录，继续访问
    }
  })
});