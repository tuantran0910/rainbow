package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tuantran0910/rainbow/internal/dtos"
	"github.com/tuantran0910/rainbow/internal/services"
	"github.com/tuantran0910/rainbow/pkg/headers"
	"github.com/tuantran0910/rainbow/pkg/utils/response"
)

type BookController struct {
	bookService services.IBookService
}

func NewBookController(bookService services.IBookService) *BookController {
	return &BookController{
		bookService: bookService,
	}
}

// GetBooks godoc
//
//	@Summary		Get Books
//	@Description	Fetch a list of books
//	@Tags			Book
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int	false	"Page number"
//	@Param			limit	query		int	false	"Number of items per page"
//	@Success		200		{object}	response.APIResponse
//	@Failure		400		{object}	response.APIResponse
//	@Failure		500		{object}	response.APIResponse
//	@Router			/books [get]
func (pc *BookController) GetBooks(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusBadRequest).
			SetMessage("Invalid pagination's parameter page, page must be a positive integer").
			SetError(err.Error()).Respond(ctx)
		return
	}

	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusBadRequest).
			SetMessage("Invalid pagination's parameter limit, limit must be a positive integer").
			SetError(err.Error()).
			Respond(ctx)
		return
	}

	reqCtx := ctx.Request.Context()
	books, pagination, err := pc.bookService.GetBooks(reqCtx, page, limit)
	if err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Failed to fetch a list of books").
			SetError(err.Error()).Respond(ctx)
		return
	}

	bookResponses := make([]*dtos.GetBookResponse, 0)
	for _, book := range books {
		bookResponses = append(bookResponses, &dtos.GetBookResponse{
			ID:            book.ID,
			SecondaryID:   book.SecondaryID,
			CategoryID:    book.CategoryID,
			SellerID:      book.SellerID,
			Name:          book.Name,
			Description:   book.Description,
			Price:         book.Price,
			OriginalPrice: book.OriginalPrice,
			RatingAverage: book.RatingAverage,
			ReviewCount:   book.ReviewCount,
			PageCount:     book.PageCount,
			SoldCount:     book.SoldCount,
			CreatedAt:     book.CreatedAt,
			UpdatedAt:     book.UpdatedAt,
			DeletedAt:     book.DeletedAt,
		})
	}
	data := &dtos.ListBooksResponse{
		Books: bookResponses,
	}

	headers := headers.NewHeaders(data, ctx)
	response.NewAPIResponse().
		SetHeaders(headers).
		SetStatusCode(http.StatusOK).
		SetPagination(pagination).
		SetMessage("Successfully retrieved books").
		SetData(data).
		Respond(ctx)
}

// GetBookById godoc
//
//	@Summary		Get Book
//	@Description	Fetch a book by its ID or secondary ID
//	@Tags			Book
//	@Accept			json
//	@Produce		json
//	@Param			id			path		string	true	"Book ID or Secondary ID"
//	@Param			secondary	query		bool	false	"Use secondary ID"
//	@Success		200			{object}	response.APIResponse
//	@Failure		404			{object}	response.APIResponse
//	@Failure		500			{object}	response.APIResponse
//	@Router			/books/{id} [get]
func (pc *BookController) GetBookById(ctx *gin.Context) {
	id := ctx.Param("id")
	isSecondary := ctx.Query("secondary") == "true"

	var bookId interface{}
	var err error
	if !isSecondary {
		bookId, err = uuid.Parse(id)
		if err != nil {
			response.
				NewAPIResponse().
				SetStatusCode(http.StatusInternalServerError).
				SetMessage("Cannot parse the ID into UUID type").
				SetError(err.Error()).Respond(ctx)
			return
		}
	} else {
		bookId = id
	}

	reqCtx := ctx.Request.Context()
	book, err := pc.bookService.GetBookById(reqCtx, bookId, isSecondary)
	if err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Failed to fetch the book").
			SetError(err.Error()).Respond(ctx)
		return
	}

	if book == nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusNotFound).
			SetMessage("Book not found").
			Respond(ctx)
		return
	}

	bookAuthors := make([]*dtos.GetAuthorResponse, 0)
	for _, author := range book.Authors {
		bookAuthors = append(bookAuthors, &dtos.GetAuthorResponse{
			ID:        author.ID,
			Name:      author.Name,
			Slug:      author.Slug,
			CreatedAt: author.CreatedAt,
			UpdatedAt: author.UpdatedAt,
			DeletedAt: author.DeletedAt,
		})
	}
	data := &dtos.GetBookResponse{
		ID:            book.ID,
		SecondaryID:   book.SecondaryID,
		CategoryID:    book.CategoryID,
		SellerID:      book.SellerID,
		Name:          book.Name,
		Description:   book.Description,
		Price:         book.Price,
		OriginalPrice: book.OriginalPrice,
		RatingAverage: book.RatingAverage,
		ReviewCount:   book.ReviewCount,
		PageCount:     book.PageCount,
		SoldCount:     book.SoldCount,
		CreatedAt:     book.CreatedAt,
		UpdatedAt:     book.UpdatedAt,
		DeletedAt:     book.DeletedAt,
		Stock: &dtos.GetInventoryResponse{
			ID:              book.Inventory.ID,
			BookID:          book.Inventory.BookID,
			Stock:           book.Inventory.Stock,
			LastRestockedAt: book.Inventory.LastRestockedAt,
			CreatedAt:       book.Inventory.CreatedAt,
			UpdatedAt:       book.Inventory.UpdatedAt,
			DeletedAt:       book.Inventory.DeletedAt,
		},
		Authors: bookAuthors,
	}

	headers := headers.NewHeaders(data, ctx)
	response.NewAPIResponse().
		SetHeaders(headers).
		SetStatusCode(http.StatusOK).
		SetMessage("Successfully retrieved book").
		SetData(data).
		Respond(ctx)
}

// CreateBook godoc
//
//	@Summary		Create Book
//	@Description	Create a book
//	@Tags			Book
//	@Accept			json
//	@Produce		json
//	@Param			req	body		dtos.CreateBookRequest	true	"Create Book Request"
//	@Success		204	{object}	response.APIResponse
//	@Failure		400	{object}	response.APIResponse
//	@Failure		500	{object}	response.APIResponse
//	@Router			/books/{id} [post]
func (pc *BookController) CreateBook(ctx *gin.Context) {
	var bookRequest dtos.CreateBookRequest
	if err := ctx.ShouldBindJSON(&bookRequest); err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusBadRequest).
			SetMessage("Invalid Body Request").
			SetError(err.Error()).
			Respond(ctx)
		return
	}

	currentUserId, ok := ctx.Get("user_id")
	if !ok {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Cannot get the current user's ID").
			Respond(ctx)
		return
	}

	reqCtx := ctx.Request.Context()
	if err := pc.bookService.CreateBook(reqCtx, bookRequest, currentUserId.(uuid.UUID)); err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusBadRequest).
			SetMessage("Cannot create book").
			SetError(err.Error()).
			Respond(ctx)
		return
	}

	headers := headers.NewHeaders(nil, ctx)
	response.NewAPIResponse().
		SetHeaders(headers).
		SetStatusCode(http.StatusCreated).
		SetMessage("Successfully created book").
		Respond(ctx)
}

// UpdateBook godoc
//
//	@Summary		Update Book
//	@Description	Update a book by its ID or secondary ID
//	@Tags			Book
//	@Accept			json
//	@Produce		json
//	@Param			id			path		string					true	"Book ID or Secondary ID"
//	@Param			secondary	query		bool					false	"Use secondary ID"
//	@Param			req			body		dtos.UpdateBookRequest	true	"Update Book Request"
//	@Success		204			{object}	response.APIResponse
//	@Failure		400			{object}	response.APIResponse
//	@Failure		500			{object}	response.APIResponse
//	@Router			/books/{id} [patch]
func (pc *BookController) UpdateBook(ctx *gin.Context) {
	id := ctx.Param("id")
	isSecondary := ctx.Query("secondary") == "true"

	var bookId interface{}
	var err error
	if !isSecondary {
		bookId, err = uuid.Parse(id)
		if err != nil {
			response.
				NewAPIResponse().
				SetStatusCode(http.StatusInternalServerError).
				SetMessage("Cannot parse the ID into UUID type").
				SetError(err.Error()).
				Respond(ctx)
			return
		}
	} else {
		bookId = id
	}

	var bookRequest dtos.UpdateBookRequest
	if err := ctx.ShouldBindJSON(&bookRequest); err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusBadRequest).
			SetMessage("Invalid Body Request").
			SetError(err.Error()).
			Respond(ctx)
		return
	}

	currentUserId, ok := ctx.Get("user_id")
	if !ok {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Cannot get the current user's ID").
			Respond(ctx)
		return
	}

	reqCtx := ctx.Request.Context()
	if err := pc.bookService.UpdateBook(reqCtx, bookId, bookRequest, currentUserId.(uuid.UUID), isSecondary); err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusBadRequest).
			SetMessage("Cannot update book").
			SetError(err.Error()).
			Respond(ctx)
		return
	}

	headers := headers.NewHeaders(nil, ctx)
	response.NewAPIResponse().
		SetHeaders(headers).
		SetStatusCode(http.StatusOK).
		SetMessage("Successfully updated book").
		Respond(ctx)
}

// DeleteBook godoc
//
//	@Summary		Delete Book
//	@Description	Delete a book by its ID
//	@Tags			Book
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Book ID"
//	@Success		204	{object}	response.APIResponse
//	@Failure		400	{object}	response.APIResponse
//	@Failure		500	{object}	response.APIResponse
//	@Router			/books/{id} [delete]
func (pc *BookController) DeleteBook(ctx *gin.Context) {
	bookId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Cannot parse the ID into UUID type").
			SetError(err.Error()).
			Respond(ctx)
		return
	}

	currentUserId, ok := ctx.Get("user_id")
	if !ok {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusInternalServerError).
			SetMessage("Cannot get the current user's ID").
			Respond(ctx)
		return
	}

	reqCtx := ctx.Request.Context()
	if err := pc.bookService.DeleteBook(reqCtx, bookId, currentUserId.(uuid.UUID)); err != nil {
		response.
			NewAPIResponse().
			SetStatusCode(http.StatusBadRequest).
			SetMessage("Failed to delete the book").
			SetError(err.Error()).
			Respond(ctx)
		return
	}

	headers := headers.NewHeaders(nil, ctx)
	response.NewAPIResponse().
		SetHeaders(headers).
		SetStatusCode(http.StatusNoContent).
		Respond(ctx)
}
