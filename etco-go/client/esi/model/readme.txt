Models are 1:1 mappings of ESI endpoint data. Endpoints with pages return a page stream, and endpoints without return a single object. For page streams, every page is returned for the query. Data transformations, at the model level, consist only of commenting out fields that aren't used by anything in the program.

Some model clients have caching-layer aliases if they're intended to be used without any transformations.
