package images

import (
    //"net/http"
)

type (
        
    // Common interface for all backends
    ImageStore interface {
        // some backends may require initialization/shutdown
        Init(map[string]string) error
        Shutdown() error
        
        // handle images api requests
        ListImages()
        GetImage()
        PutImage()
        DeleteImage()
    }
    
)