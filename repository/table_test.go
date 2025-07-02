package repository

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDynamoDBRepo_CreateTable(t *testing.T) {
	ctx := context.Background()
	tableName := "test-table"

	t.Run("should create table and wait for it to become active", func(t *testing.T) {
		// 1. Setup
		mockClient := new(MockDynamoDBClient)
		repo, err := NewDynamoDBRepo(mockClient, tableName)
		require.NoError(t, err)

		// 2. Expectations
		// Esperamos a chamada inicial para criar a tabela. Retornamos sucesso.
		mockClient.On("CreateTable", mock.Anything, mock.AnythingOfType("*dynamodb.CreateTableInput")).
			Return(&dynamodb.CreateTableOutput{}, nil).Once()
		// Esperamos a chamada para DescribeTable, que é feita pelo Waiter.
		// Simulamos que a tabela agora está ativa, o que fará o waiter parar e retornar sucesso.
		mockClient.On("DescribeTable", mock.Anything, mock.AnythingOfType("*dynamodb.DescribeTableInput")).
			Return(&dynamodb.DescribeTableOutput{
				Table: &types.TableDescription{
					TableStatus: types.TableStatusActive,
				},
			}, nil).
			Once()

		// 3. Execute
		// A função de "action" pode ser vazia, pois não estamos testando o ticker.
		err = repo.CreateTable(ctx, func() {})

		// 4. Assert
		assert.NoError(t, err)
		mockClient.AssertExpectations(t)
	})

	// ... outros cenários de teste, como a tabela já existir ...
}
