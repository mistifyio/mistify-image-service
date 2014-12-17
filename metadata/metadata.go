package metadata

type (
    
    // Common interface for metadata stores
    MetadataStore interface {
        // some backends may require initialization/shutdown
        Init(map[string]string) error
        Shutdown() error
        
        // metadata api request handlers
        
    }
    
)