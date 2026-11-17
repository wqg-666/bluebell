
<template>
  <div class="content">
    <div class="left">
      <div class="post-name">我好想写点什么</div>
      <div class="post-type">
        <input type="text" class="post-type-value" placeholder="选择一个频道" :value="selectCommunity.name" readonly @click="showCommunity()"/>
        <ul class="post-type-options" v-show="showCommunityList">
          <li class="post-type-cell"
            v-for="(community, index) in communityList"
            :key="community.id"
            @click="selected(index)"
          >
            {{community.name}}
          </li>
        </ul>
        <i class="p-icon"></i>
      </div>
      <div class="post-content">
        <ul class="cat">
          <li class="cat-item active">
            <i class="iconfont icon-edit"></i>post
          </li>
          <li class="cat-item">
            <i class="iconfont icon-image"></i>image/video
          </li>
        </ul>
        <div class="post-sub-container">
          <div class="post-sub-header">
            <textarea class="post-title" cols="30" rows="10" v-model="title" maxlength="300" placeholder="标题"></textarea>
            <span class="textarea-num">{{ title.length }}/300</span>
          </div>
          <!---此处放置富文本--->
          <div class="post-text-con">
            <textarea
              class="post-content-t"
              cols="30"
              rows="10"
              v-model="content"
              placeholder="内容"
            ></textarea>
          </div>
        </div>
        <p v-if="formError" class="form-error" role="alert">{{ formError }}</p>
        <div class="post-footer">
          <div class="btns">
            <button class="btn">取消</button>
            <button class="btn" :disabled="isSubmitting" @click="submit()">{{ isSubmitting ? "发表中…" : "发表" }}</button>
          </div>
        </div>
      </div>
    </div>
    <div class="right">
      <div class="post-rank">
        <h5 class="p-r-title">
          <i class="p-r-icon"></i>发帖规范
        </h5>
        <ul class="p-r-content">
          <li class="p-r-item">1.网络不是法外之地</li>
          <li class="p-r-item">2.网络不是法外之地</li>
          <li class="p-r-item">3.网络不是法外之地</li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "Publish",
  data() {
    return {
      title: "",
      content: "",
      showCommunityList: false,
      selectCommunity: {},
      communityList: [],
      formError: "",
      isSubmitting: false
    };
  },
  methods: {
    async submit() {
      if (this.isSubmitting) return;
      this.formError = "";
      // 频道接口的 id 是字符串，创建帖子接口的 community_id 是 int64。
      const communityID = Number(this.selectCommunity.id);
      if (!Number.isSafeInteger(communityID) || communityID <= 0) {
        this.formError = "请选择一个有效的频道";
        return;
      }
      if (!this.title.trim() || !this.content.trim()) {
        this.formError = "请填写标题和内容";
        return;
      }
      if (this.title.length > 300) {
        this.formError = "标题不能超过 300 个字符";
        return;
      }

      this.isSubmitting = true;
      try {
        const response = await this.$axios({
          method: "post",
          url: "/post",
          data: {
            title: this.title.trim(),
            content: this.content,
            community_id: communityID
          }
        });
        if (response.code == 1000) {
          this.$router.push({ path: "/" });
        } else {
          this.formError = typeof response.msg === "string" ? response.msg : "发表失败，请稍后重试";
        }
      } catch (error) {
        this.formError = "发表失败，请检查网络连接后重试";
      } finally {
        this.isSubmitting = false;
      }
    },
    getCommunityList() {
      this.$axios({
        method: "get",
        url: "/community"
      })
        .then(response => {
          if (response.code == 1000) {
            this.communityList = response.data || [];
          } else {
            this.formError = typeof response.msg === "string" ? response.msg : "频道加载失败，请刷新后重试";
          }
        })
        .catch(() => {
          this.formError = "频道加载失败，请检查网络连接后刷新";
        });
    },
    showCommunity(){
      this.showCommunityList = !this.showCommunityList;
    },
    selected(index) {
      this.selectCommunity = this.communityList[index];
      this.showCommunityList = false;
    }
  },
  mounted: function() {
    this.getCommunityList();
  }
};
</script>
<style lang="less" scoped>
.content {
  max-width: 100%;
  box-sizing: border-box;
  display: flex;
  flex-direction: row;
  justify-content: center;
  margin: 0 auto;
  padding: 20px 24px;
  margin-top: 48px;
  .left {
    flex-grow: 1;
    max-width: 740px;
    min-width: 0;
    word-break: break-word;
    flex: 1;
    margin: 32px;
    margin-right: 12px;
    padding-bottom: 30px;
    position: relative;
    .post-name {
      padding: 4px;
      margin: 16px 0;
      border-bottom: solid 1px #edeff1;
      display: -webkit-flex;
      display: flex;
      justify-content: space-between;
      .p-btn {
        font-size: 12px;
        font-weight: 700;
        letter-spacing: 0.5px;
        line-height: 24px;
        text-transform: uppercase;
        border: none;
        padding: 0;
        margin-left: 10px;
        color: #0079d3;
      }
      .p-num {
        font-size: 12px;
        font-weight: 400;
        line-height: 16px;
        background: #878a8c;
        border-radius: 2px;
        color: #ffffff;
        margin-left: 4px;
        padding: 1px 3px;
      }
    }
    .post-type {
      position: relative;
      box-sizing: border-box;
      width: 300px;
      max-width: 100%;
      height: 40px;
      border-radius: 4px;
      transition: box-shadow 0.2s ease;
      box-shadow: 0 0 0 0 #ffffff;
      border: 1px solid #edeff1;
      background-color: #ffffff;
      padding-left: 10px;
      position: relative;
      .post-type-value {
        font-size: 14px;
        font-weight: 500;
        line-height: 40px;
        width: 100%;
        vertical-align: middle;
        color: #1c1c1c;
        background-color: transparent;
        cursor: pointer;
      }
      .post-type-options {
        position: absolute;
        width: 100%;
        background-color: white;
        left: 0;
        z-index: 1;
        border-radius: 4px;
        .post-type-cell {
          margin: 14px 8px 5px;
          font-size: 14px;
          list-style: none;
          border-bottom: 1px solid #edeff1;
          padding-bottom: 8px;
          color: #1c1c1c;
          cursor: pointer;
        }
      }
      .p-icon {
        width: 0;
        height: 0;
        border-top: 5px solid #878a8c;
        border-right: 5px solid transparent;
        border-bottom: 5px solid transparent;
        border-left: 5px solid transparent;
        margin-left: 10px;
        position: absolute;
        top: 50%;
        right: 10px;
        cursor: pointer;
      }
    }
    .post-content {
      background-color: #ffffff;
      margin: 10px 0;
      padding-bottom: 15px;
      border-radius: 5px;
      .cat {
        display: flex;
        display: -webkit-flex;
        justify-content: space-between;
        align-items: center;
        width: 100%;
        height: 53px;
        .cat-item {
          padding: 10px 0;
          width: 50%;
          height: 40px;
          line-height: 40px;
          text-align: center;
          list-style: none;
          border-bottom: 1px solid #edeff1;
          border-right: 1px solid #edeff1;
          color: #878a8c;
          .iconfont {
            margin-right: 4px;
          }
        }
        .active {
          color: #0079d3;
          font-weight: bolder;
          background: none;
        }
      }
      .post-sub-container {
        padding: 16px;
        .post-sub-header {
          position: relative;
          .post-title {
            resize: none;
            box-sizing: border-box;
            overflow: hidden;
            display: block;
            width: 100%;
            height: 40px;
            padding: 0 70px 0 10px;
            outline: none;
            border: 1px solid #edeff1;
            border-radius: 4px;
            color: #1c1c1c;
            font-size: 14px;
            font-weight: 400;
            line-height: 40px;
          }
          .textarea-num {
            font-size: 10px;
            font-weight: 700;
            letter-spacing: 0.5px;
            line-height: 12px;
            text-transform: uppercase;
            bottom: 12px;
            color: #878a8c;
            pointer-events: none;
            position: absolute;
            right: 12px;
          }
        }
        .post-text-con {
          width: 100%;
          height: 200px;
          border: 1px solid #edeff1;
          margin-top: 20px;
          .post-content-t {
            resize: none;
            box-sizing: border-box;
            overflow: hidden;
            display: block;
            width: 100%;
            height: 200px;
            padding: 12px 8px;
            outline: none;
            border: 1px solid #edeff1;
            border-radius: 4px;
            color: #1c1c1c;
            font-size: 14px;
            font-weight: 400;
            line-height: 21px;
          }
        }
      }
      .post-footer {
        display: flex;
        display: -webkit-flex;
        margin: 0 16px;
        justify-content: flex-end;
        .sign {
          display: flex;
          display: -webkit-flex;
          .sign-item {
            list-style: none;
            padding: 5px 8px;
            border: 1px solid #edeff1;
            margin-right: 10px;
            color: #878a8c;
            font-size: 12px;
            font-weight: 700;
          }
        }
        .btns {
          .btn {
            border: 1px solid transparent;
            border-radius: 4px;
            box-sizing: border-box;
            text-align: center;
            text-decoration: none;
            font-size: 12px;
            font-weight: 700;
            letter-spacing: 0.5px;
            line-height: 24px;
            text-transform: uppercase;
            padding: 3px 16px;
            background: #0079d3;
            color: #ffffff;
            margin-left: 8px;
            cursor: pointer;
            &:disabled {
              opacity: 0.6;
              cursor: wait;
            }
          }
        }
      }
      .alias {
        background-color: #f6f7f8;
        border-radius: 0 0 6px 6px;
        border-top: solid 1px #edeff1;
        display: -ms-flexbox;
        display: flex;
        -ms-flex-flow: column;
        flex-flow: column;
        padding: 8px 16px 21px;
        position: relative;
        .send-post {
          font-size: 14px;
          font-weight: 500;
          line-height: 18px;
          color: #1c1c1c;
          margin-right: 4px;
        }
        .connect {
          font-size: 14px;
          font-weight: 500;
          line-height: 18px;
          color: #0079d3;
          display: block;
          margin-right: 4px;
          margin-top: 10px;
        }
      }
    }
  }
  .right {
    flex-grow: 0;
    flex-shrink: 0;
    width: 312px;
    margin-top: 62px;
    .post-rank {
      background-color: #ffffff;
      border-radius: 4px;
      margin-top: 15px;
      padding: 12px;
      .p-r-title {
        display: flex;
        display: -webkit-flex;
        align-items: center;
        .p-r-icon {
          width: 40px;
          height: 40px;
          background: url("../assets/images/avatar.png") no-repeat;
          background-size: cover;
          margin-right: 10px;
        }
        font-size: 16px;
        font-weight: 500;
        line-height: 20px;
        -ms-flex-align: center;
        align-items: center;
        border-bottom: 1px solid #edeff1;
        color: #1c1c1c;
        padding-bottom: 10px;
        // display: -ms-flexbox;
        // display: flex;
      }
      .p-r-content {
        display: flex;
        display: -webkit-flex;
        flex-direction: column;
        .p-r-item {
          list-style: none;
          border-bottom: 1px solid #edeff1;
          color: #1c1c1c;
          padding: 10px 5px;
        }
      }
    }
  }
}
.form-error {
  margin: 0 16px 12px;
  color: #c62828;
  font-size: 14px;
}
@media (max-width: 768px) {
  .content {
    flex-direction: column;
    padding: 16px 12px;
    .left {
      flex: none;
      width: 100%;
      margin: 0;
      padding-bottom: 0;
    }
    .right {
      width: 100%;
      margin-top: 0;
    }
  }
}
</style>
