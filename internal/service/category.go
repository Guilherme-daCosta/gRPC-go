package service

import (
	"context"
	"io"

	"github.com/Guilherme-daCosta/gRPC-go/internal/database"
	"github.com/Guilherme-daCosta/gRPC-go/internal/pb"
)

type CategoryService struct {
	pb.UnimplementedCategoryServiceServer
	CategoryDB database.Category
}

func NewCategoryService(categoryDB database.Category) *CategoryService {
	return &CategoryService{
		CategoryDB: categoryDB,
	}
}

func (c *CategoryService) CreateCategory(ctx context.Context, in *pb.CreateCategoryRequest) (*pb.Category, error) {
	category, err := c.CategoryDB.Add(in.Name, in.Description)
	if err != nil {
		return nil, err
	}

	categoryResponse := &pb.Category{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}

	return categoryResponse, nil
}

func (c *CategoryService) ListCategories(ctx context.Context, in *pb.Empty) (*pb.CategoryList, error) {
	categories, err := c.CategoryDB.FindAll()
	if err != nil {
		return nil, err
	}

	categoryList := &pb.CategoryList{}

	for _, category := range categories {
		categoryList.Categories = append(categoryList.Categories, &pb.Category{
			Id:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		})
	}

	return categoryList, nil
}

func (c *CategoryService) GetCategory(ctx context.Context, in *pb.CategoryGetRequest) (*pb.Category, error) {
	category, err := c.CategoryDB.FindByID(in.Id)
	if err != nil {
		return nil, err
	}

	categoryResponse := &pb.Category{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}

	return categoryResponse, nil
}

func (c *CategoryService) CreateCategoryStream(stream pb.CategoryService_CreateCategoryStreamServer) error {
	categories := &pb.CategoryList{}
	
	for {
		category, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(categories)
		}
		if err != nil {
			return err
		}

		categoryResult, err := c.CategoryDB.Add(category.Name, category.Description)
		if err != nil {
			return err
		}
		categories.Categories = append(categories.Categories, &pb.Category{
			Id:          categoryResult.ID,
			Name:        categoryResult.Name,
			Description: categoryResult.Description,
		})
	}
}

func (c *CategoryService) CreateCategoryStreamBidirectional(stream pb.CategoryService_CreateCategoryStreamBidirectionalServer) error {
	for {
		category, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		categoryResult, err := c.CategoryDB.Add(category.Name, category.Description)
		if err != nil {
			return err
		}
		
		categoryResponse := &pb.Category{
			Id:          categoryResult.ID,
			Name:        categoryResult.Name,
			Description: categoryResult.Description,
		}
		if err := stream.Send(categoryResponse); err != nil {
			return err
		}
	}
}