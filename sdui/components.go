package sdui

// The standard-component type names — the vocabulary a producer emits and a
// client renders from the shared definition library (../definitions). Screens
// are authored with the declarative ui layer (github.com/mosaic-media/contracts/ui),
// which builds these node types; these constants name them for any code that
// needs to match on a type.

// Node type names.
const (
	TypeScreen          = "Screen"
	TypeSection         = "Section"
	TypeStack           = "Stack"
	TypeGrid            = "Grid"
	TypeCarousel        = "Carousel"
	TypeDivider         = "Divider"
	TypePosterCard      = "PosterCard"
	TypeHeroBanner      = "HeroBanner"
	TypeDetailHeader    = "DetailHeader"
	TypeEpisodeRow      = "EpisodeRow"
	TypeSeasonSelector  = "SeasonSelector"
	TypeRelatedRail     = "RelatedRail"
	TypeSourcePicker    = "SourcePicker"
	TypePlaybackBar     = "PlaybackBar"
	TypePersonChip      = "PersonChip"
	TypeGenreTag        = "GenreTag"
	TypeButton          = "Button"
	TypeIconButton      = "IconButton"
	TypeBadge           = "Badge"
	TypeBanner          = "Banner"
	TypeStatusIndicator = "StatusIndicator"
	TypeEmptyState      = "EmptyState"
	TypeSearchBar       = "SearchBar"
	TypeTextField       = "TextField"
	TypeToggle          = "Toggle"
	TypeSelect          = "Select"
	TypeSlider          = "Slider"
	TypeProgressBar     = "ProgressBar"
	TypePagination      = "Pagination"
)
