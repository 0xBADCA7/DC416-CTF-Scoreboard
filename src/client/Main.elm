module Main exposing (main)

import Html exposing (..)


-- MAIN


main : Platform.Program Basics.Never Model Msg
main =
    Html.program
        { init = init
        , update = update
        , view = view
        , subscriptions = subscriptions
        }



-- MODEL


type alias Date =
    String


type ViewMode
    = ScoreboardView
    | FormView
    | MessagesView


type alias Team =
    { name : String
    , rank : Int
    , score : Int
    , lastSubmitted : Date
    }


type alias Scoreboard =
    List Team


type Message
    = Message Date String


type alias Model =
    { mode : ViewMode
    , scoreboard : Scoreboard
    , messages : List Message
    }


init : ( Model, Cmd Msg )
init =
    ( Model ScoreboardView emptyScoreboard [], Cmd.none )


emptyScoreboard : Scoreboard
emptyScoreboard =
    []



-- UPDATE


type Msg
    = SwitchMode ViewMode


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        SwitchMode mode ->
            ( { model | mode = mode }, Cmd.none )



-- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none



-- VIEW


view : Model -> Html Msg
view model =
    div []
        [ p [] [ text "Hello, World!" ]
        ]
