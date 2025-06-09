package handlers

import (
	"fmt"
	"net/http"

	"github.com/kevinjimenez96/chirpy/internal/types"
)

func MetricsHandler(w http.ResponseWriter, r *http.Request, cfg *types.ApiConfig) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(([]byte)(fmt.Sprintf(metricsBody, cfg.FileserverHits.Load())))
}

var metricsBody string = `
<html>

<body>
	<h1>Welcome, Chirpy Admin</h1>
	<p>Chirpy has been visited %d times!</p>
</body>

</html>
`
