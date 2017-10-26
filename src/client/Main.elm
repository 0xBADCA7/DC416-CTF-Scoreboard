module Main exposing (main)

-- Standard imports

import Http
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (..)


-- Local imports

import Mode.Scoreboard as Scoreboard exposing (Scoreboard(..))
import Mode.Message as Message exposing (Message)
import Mode.SubmitForm as SubmitForm exposing (SubmitResponse)


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


type Notification
    = Error String
    | Success String


type ViewMode
    = ScoreboardView
    | SubmitForm
    | MessagesView


type Msg
    = SwitchMode ViewMode
    | ScoreboardRetrieved (Result Http.Error Scoreboard)
    | FlagSubmitted (Result Http.Error SubmitResponse)
    | MessagesRetrieved (Result Http.Error (List Message))
    | GotInput SubmitForm.Input


type alias Model =
    { scoreboard : Scoreboard
    , submitTokenInput : String
    , flagInput : String
    , messages : List Message
    , mode : ViewMode
    , notification : Maybe Notification
    }


init : ( Model, Cmd Msg )
init =
    ( Model (Scoreboard []) "" "" [] ScoreboardView Nothing, Scoreboard.query ScoreboardRetrieved )



-- UPDATE


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        SwitchMode ScoreboardView ->
            ( { model | mode = ScoreboardView }, Scoreboard.query ScoreboardRetrieved )

        SwitchMode MessagesView ->
            ( { model | mode = MessagesView }, Message.query MessagesRetrieved )

        SwitchMode mode ->
            ( { model | mode = mode }, Cmd.none )

        ScoreboardRetrieved (Ok scoreboard) ->
            ( { model | scoreboard = scoreboard }, Cmd.none )

        ScoreboardRetrieved _ ->
            ( model, Cmd.none )

        MessagesRetrieved (Ok messages) ->
            ( { model | messages = messages }, Cmd.none )

        MessagesRetrieved (Err err) ->
            let
                _ =
                    Debug.log "ERROR" err
            in
                ( model, Cmd.none )

        GotInput (SubmitForm.SubmissionToken input) ->
            ( { model | submitTokenInput = input }, Cmd.none )

        GotInput (SubmitForm.Flag input) ->
            ( { model | flagInput = input }, Cmd.none )

        GotInput SubmitForm.Submit ->
            ( model, SubmitForm.mutation model.submitTokenInput model.flagInput FlagSubmitted )

        FlagSubmitted (Err _) ->
            ( { model | notification = Just <| Error "Error submitting flag. Try again later." }, Cmd.none )

        FlagSubmitted (Ok { correct, newScore, scoreboard }) ->
            let
                notification =
                    if correct then
                        Success ("Congratulations! Your team's score is now" ++ (toString newScore) ++ "!")
                    else
                        Error "Your submission was incorrect."
            in
                ( { model | notification = Just notification, scoreboard = scoreboard }, Cmd.none )



-- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none



-- VIEW


view : Model -> Html Msg
view model =
    div []
        [ viewNav model
        , viewMode model
        ]


viewNav : Model -> Html Msg
viewNav model =
    let
        navLinks =
            case model.mode of
                ScoreboardView ->
                    [ li [ class "active" ] [ a [ href "#", onClick (SwitchMode ScoreboardView) ] [ text "Scoreboard" ] ]
                    , li [] [ a [ href "#", onClick (SwitchMode MessagesView) ] [ text "Messages" ] ]
                    , li [] [ a [ href "#", onClick (SwitchMode SubmitForm) ] [ text "Submit" ] ]
                    ]

                MessagesView ->
                    [ li [] [ a [ href "#", onClick (SwitchMode ScoreboardView) ] [ text "Scoreboard" ] ]
                    , li [ class "active" ] [ a [ href "#", onClick (SwitchMode MessagesView) ] [ text "Messages" ] ]
                    , li [] [ a [ href "#", onClick (SwitchMode SubmitForm) ] [ text "Submit" ] ]
                    ]

                SubmitForm ->
                    [ li [] [ a [ href "#", onClick (SwitchMode ScoreboardView) ] [ text "Scoreboard" ] ]
                    , li [] [ a [ href "#", onClick (SwitchMode MessagesView) ] [ text "Messages" ] ]
                    , li [ class "active" ] [ a [ href "#", onClick (SwitchMode SubmitForm) ] [ text "Submit" ] ]
                    ]
    in
        nav []
            [ div [ class "mainContent nav-wrapper" ]
                [ a [ href "#", class "brand-logo" ] [ text "Scoreboard" ]
                , ul [ id "nav-mobile", class "right" ] navLinks
                ]
            ]


viewMode : Model -> Html Msg
viewMode model =
    let
        viewContent =
            case model.mode of
                ScoreboardView ->
                    [ span [ class "card-title gray-text text-darken-4" ] [ text "Scoreboard" ]
                    , Scoreboard.view model.scoreboard
                    ]

                SubmitForm ->
                    [ span [ class "card-title gray-text text-darken-4" ] [ text "Submit a flag" ]
                    , SubmitForm.view GotInput
                    ]

                MessagesView ->
                    [ span [ class "card-title gray-text text-darken-4" ] [ text "Messages" ]
                    , Message.view model.messages
                    ]
    in
        div [ class "mainContent" ]
            [ div [ id "content", class "card large waves-effect waves-light" ]
                [ div [ class "card-content" ] viewContent
                ]
            ]
