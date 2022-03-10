package newsletter

import (
	"context"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/sankethkini/NewsLetter-Backend/internal/kproducer"
	v1 "github.com/sankethkini/NewsLetter-Backend/proto/newsletterpb/v1"
	"github.com/stretchr/testify/require"
)

func TestCreateNewsLetterSuccess(t *testing.T) {
	ctx := context.Background()
	t.Logf("Name:%s", t.Name())
	ctrl := gomock.NewController(t)
	mockDB := NewMockDB(ctrl)
	mockKafka := kproducer.NewMockProducer(ctrl)
	svc := NewNewsService(mockDB, mockKafka)
	mockDB.EXPECT().addNewsLetter(gomock.Any(), gomock.Any()).Return(&NewsLetterModel{Title: "some title", NewsLetterID: "some123", Body: "kjs sfnsakjnf askfjas"}, nil).Times(1)
	resp, err := svc.CreateNewsLetter(ctx, &v1.CreateNewsLetterRequest{Title: "some title", Body: "kjs sfnsakjnf askfjas"})
	require.NotNil(t, resp)
	require.Nil(t, err)
}

func TestAddSchemeToNewsSuccess(t *testing.T) {
	ctx := context.Background()
	t.Logf("Name:%s", t.Name())
	ctrl := gomock.NewController(t)
	mockDB := NewMockDB(ctrl)
	mockKafka := kproducer.NewMockProducer(ctrl)
	svc := NewNewsService(mockDB, mockKafka)
	mockDB.EXPECT().addSchemeToNews(gomock.Any(), gomock.Any()).Return(&NewsSchemes{NewsLetterID: "someid123", SchemeID: "someid123"}, nil).Times(1)
	mockDB.EXPECT().getNewsLetter(gomock.Any(), gomock.Any()).Return(&NewsLetterModel{NewsLetterID: "someid123", Title: "some title", Body: "some body"}, nil).Times(1)
	mockKafka.EXPECT().Produce(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
	resp, err := svc.AddSchemeToNews(ctx, &v1.NewsScheme{NewsLetterId: "someid123", SchemeId: "someid123"})
	require.NotNil(t, resp)
	require.Nil(t, err)
}
