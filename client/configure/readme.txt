RawData can be broken down to this: map[ShopLocationKey]map[TypeId]Price
Web API Data can be broken down to this: map[TypeId]map[ShopLocationKey]Price

In reality, it's moreso that the raw data has "ShopLocations" which each contains
    a "map[TypeId]Price." But, the above is easier to understand.

I've worked with ChatGPT to figure out how to name these data structures, and signify
    their differences and similarities.

"""
Those names seem clear and more intuitive! Here's a brief breakdown of what I infer
    from the naming:

    1. **ShopLocation**: This sounds like a parent entity, probably representing a
        particular physical or virtual shop/store. It gives the idea of a centralized
        or main entity.
    2. **ShopLocationTypeMap**: This name gives a clear indication that it's mapping
        types (probably product types or item types) to their respective prices for
        a particular shop location. The use of "child" in your description implies
        this map might be contained within or associated with a specific `ShopLocation`.

    3. **ShopLocationTypeMapsBuilder**: The addition of "Builder" implies that this
        is an interface or utility to help construct or manage the type maps. This
        makes sense for a web context, where you'd be manipulating or constructing
        these maps based on user input or configuration.
    4. **ShopLocationTypeBundle**: The term "Bundle" implies a collection or grouping,
        which makes sense if it contains various price keys (or location keys) and
        their associated prices. 

These names are more descriptive than the original ones and seem to give a clearer
    picture of their purpose and their relationships. The naming also aligns well
    with the concept of locations, types, and prices, which is the core functionality
    of your data structures. Good job on refining them!
"""
