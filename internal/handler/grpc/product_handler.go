package grpc

import (
	"context"

	"grpc-exmpl/internal/model"
	"grpc-exmpl/internal/service"
	pb "grpc-exmpl/proto/product"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductHandler struct {
	pb.UnimplementedProductServiceServer
	service service.ProductService
}

func NewProductHandler(svc service.ProductService) *ProductHandler {
	return &ProductHandler{service: svc}
}

// CreateProduct handles gRPC request to create a new product
func (h *ProductHandler) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	productReq := &model.CreateProductRequest{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       int(req.Stock),
		UserID:      req.UserId,
	}

	product, err := h.service.CreateProduct(productReq)
	if err != nil {
		return &pb.CreateProductResponse{Success: false, Message: err.Error()}, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.CreateProductResponse{
		Success: true,
		Message: "Product created successfully",
		Product: convertProductToProto(product),
	}, nil
}

// GetProduct handles gRPC request to get a product by ID
func (h *ProductHandler) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	product, err := h.service.GetProductByID(req.Id)
	if err != nil {
		return &pb.GetProductResponse{Success: false, Message: err.Error()}, status.Error(codes.NotFound, err.Error())
	}

	return &pb.GetProductResponse{
		Success: true,
		Message: "OK",
		Product: convertProductToProto(product),
	}, nil
}

// ListProducts handles gRPC request to list products by user
func (h *ProductHandler) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	products, err := h.service.ListProductsByUser(req.UserId)
	if err != nil {
		return &pb.ListProductsResponse{Success: false, Message: err.Error()}, status.Error(codes.Internal, err.Error())
	}

	var productProtos []*pb.ProductData
	for _, p := range products {
		productProtos = append(productProtos, convertProductToProto(p))
	}

	return &pb.ListProductsResponse{
		Success:  true,
		Message:  "OK",
		Products: productProtos,
	}, nil
}

// UpdateProduct handles gRPC request to update product
func (h *ProductHandler) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	updReq := &model.UpdateProductRequest{
		ID:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       int(req.Stock),
	}

	product, err := h.service.UpdateProduct(updReq)
	if err != nil {
		return &pb.UpdateProductResponse{Success: false, Message: err.Error()}, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.UpdateProductResponse{
		Success: true,
		Message: "Product updated successfully",
		Product: convertProductToProto(product),
	}, nil
}

// DeleteProduct handles gRPC request to delete product
func (h *ProductHandler) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	if err := h.service.DeleteProduct(req.Id); err != nil {
		return &pb.DeleteProductResponse{Success: false, Message: err.Error()}, status.Error(codes.NotFound, err.Error())
	}

	return &pb.DeleteProductResponse{
		Success: true,
		Message: "Product deleted",
	}, nil
}

// convertProductToProto maps internal Product model to gRPC proto message
func convertProductToProto(p *model.Product) *pb.ProductData {
	return &pb.ProductData{
		Id:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       int32(p.Stock),
		UserId:      p.UserID,
		CreatedAt:   p.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   p.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
