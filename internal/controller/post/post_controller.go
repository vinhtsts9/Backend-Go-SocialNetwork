package post

import (
	"fmt"
	"go-ecommerce-backend-api/m/v2/global"
	model "go-ecommerce-backend-api/m/v2/internal/models"
	"go-ecommerce-backend-api/m/v2/internal/service"
	"go-ecommerce-backend-api/m/v2/package/utils/auth"
	"go-ecommerce-backend-api/m/v2/response"
	"go-ecommerce-backend-api/m/v2/worker"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

// management post
var Post = new(cPost)

type cPost struct{}

// CreatePost
// @Summary      Create a new post
// @Description  Create a new post for the user
// @Tags         post management
// @Accept       json
// @Produce      json
// @Param        payload body model.CreatePostInput true "Post Payload"
// @Success      201  {object}  response.ResponseData
// @Failure      400  {object}  response.ErrorResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router       /post/create [post]
func (c *cPost) CreatePost(ctx *gin.Context) {
	fmt.Println("DEBUG: CreatePost called")

	// Phân tích form-data
	err := ctx.Request.ParseMultipartForm(10 << 20) // Giới hạn 10MB
	if err != nil {
		fmt.Println("DEBUG: Multipart Form Parse Error:", err)
		ctx.JSON(400, gin.H{"error": "Invalid form data"})
		return
	}
	fmt.Println("DEBUG: Multipart Form Parsed Successfully")

	userNicknameStr := ctx.PostForm("user_nickname")
	title := ctx.DefaultPostForm("title", "")
	isPublishedStr := ctx.PostForm("is_published")
	metadata := ctx.DefaultPostForm("metadata", "")

	userInfo := auth.GetUserInfoFromContext(ctx)

	// Chuyển đổi is_published sang bool
	isPublished, err := strconv.ParseBool(isPublishedStr)
	if err != nil {
		fmt.Println("DEBUG: error convert isPublished", err)
		ctx.JSON(400, gin.H{"error": fmt.Sprintf("error convert isPublished: %s", err.Error())})
		return
	}

	// Nhận các file ảnh
	files := ctx.Request.MultipartForm.File["image_paths"]
	fmt.Printf("DEBUG: number of files to upload: %d\n", len(files))

	var imageUrls []string
	for i, file := range files {

		fmt.Printf("DEBUG: Processing file %d: %s\n", i, file.Filename)

		fileContent, err := file.Open()
		if err != nil {
			fmt.Println("DEBUG: Error opening file:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error opening file: %s", err.Error())})
			return
		}
		defer fileContent.Close()

		// Upload ảnh lên Cloudinary
		uploadResp, err := global.Cloudinary.UploadImageToCloudinaryFromReader(fileContent, "uploads")
		if err != nil {
			fmt.Println("DEBUG: Cloudinary upload failed:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Cloudinary upload failed: %s", err.Error())})
			return
		}

		fmt.Printf("DEBUG: Uploaded file %d URL: %s\n", i, uploadResp)
		// Thêm URL của ảnh vào danh sách
		imageUrls = append(imageUrls, uploadResp)
	}

	// In ra danh sách ảnh để kiểm tra
	fmt.Println("DEBUG: imageUrls:", imageUrls)

	// Tạo payload cho bài viết
	post := model.CreatePostInput{
		UserNickname: userNicknameStr,
		Title:        title,
		ImagePaths:   imageUrls,
		IsPublished:  isPublished,
		Metadata:     metadata,
		UserId:       uint64(userInfo.UserID),
	}
	fmt.Println("DEBUG: post payload:", post)

	// Gọi phương thức DistributeTaskPost
	opts := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(10 * time.Second),
		asynq.Queue(worker.QueueCritical),
	}
	distributor := worker.NewRedisTaskDistributor(global.RedisOpt)

	dataRs, err := distributor.DistributeTaskPost(ctx.Request.Context(), &post, opts...)
	if err != nil {
		fmt.Println("DEBUG: DistributeTaskPost failed:", err)
		ctx.JSON(500, gin.H{"error": fmt.Sprintf("Distribute post failed: %s", err.Error())})
		return
	}
	fmt.Println("DEBUG: DistributeTaskPost success:", dataRs)

	// Trả về kết quả thành công
	response.SuccessResponse(ctx, 200, dataRs)
}

// UpdatePost
// @Summary      Update a post
// @Description  Update a post by its ID
// @Tags         post management
// @Accept       json
// @Produce      json
// @Param        id   path     int  true  "Post ID"
// @Param        payload body model.UpdatePostInput true "Updated Post Data"
// @Success      200  {object}  response.ResponseData
// @Failure      400  {object}  response.ErrorResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router       /post/{id} [patch]
func (c *cPost) UpdatePost(ctx *gin.Context) {
	id := ctx.Param("id")
	var params model.UpdatePostInput
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, err.Error())
		return
	}
	codeRs, dataRs, err := service.NewPost().UpdatePost(ctx, id, &params)
	if err != nil {
		global.Logger.Error("Error updating post", zap.Error(err))
		response.ErrorResponse(ctx, response.ErrCodeInternal, err.Error())
		return
	}
	response.SuccessResponse(ctx, codeRs, dataRs)
}

// DeletePost
// @Summary      Delete a post
// @Description  Delete a post by its ID
// @Tags         post management
// @Accept       json
// @Produce      json
// @Param        id   path     int  true  "Post ID"
// @Success      204  {object}  response.ResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router       /post/{id} [delete]
func (c *cPost) DeletePost(ctx *gin.Context) {
	id := ctx.Param("id")
	codeRs, err := service.NewPost().DeletePost(ctx, id)
	if err != nil {
		global.Logger.Error("Error deleting post", zap.Error(err))
		response.ErrorResponse(ctx, response.ErrCodeInternal, err.Error())
		return
	}
	response.SuccessResponse(ctx, codeRs, nil)
}
